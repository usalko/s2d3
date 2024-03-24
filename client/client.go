package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/usalko/s2d3/utils"
	"golang.org/x/net/proxy"
)

type TraceLevel int8

const (
	NoTrace      TraceLevel = 0
	TraceHeaders TraceLevel = 1
	TraceAll     TraceLevel = 2
)

type Client struct {
	AccessKeyId     string
	SecretAccessKey string

	Token       string
	Region      string
	Bucket      string
	Domain      string
	Protocol    string
	SOCKS5Proxy string

	SignatureVersion int

	CACertificates     []string
	SkipSystemCAs      bool
	InsecureSkipVerify bool

	UsePathBuckets bool

	httpClient *http.Client

	traceWriter io.Writer
	traceLevel  TraceLevel
}

func NewClient(client *Client) (*Client, error) {
	var (
		roots *x509.CertPool
		err   error
	)

	if client.SignatureVersion == 0 {
		client.SignatureVersion = 4
	}

	if trace := os.Getenv("S3_TRACE_LEVEL"); trace != "" {
		switch strings.ToLower(trace) {
		case "1":
			client.Trace(os.Stderr, TraceHeaders)
		case "2":
			client.Trace(os.Stderr, TraceAll)
		default:
			client.Trace(os.Stderr, NoTrace)
		}
	}

	if !client.SkipSystemCAs {
		roots, err = x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve system root certificate authorities: %s", err)
		}
	} else {
		roots = x509.NewCertPool()
	}

	for _, ca := range client.CACertificates {
		if ok := roots.AppendCertsFromPEM([]byte(ca)); !ok {
			return nil, fmt.Errorf("unable to append CA certificate")
		}
	}

	dialContext := http.DefaultTransport.(*http.Transport).DialContext
	if client.SOCKS5Proxy != "" {
		dialer, err := proxy.SOCKS5("tcp", client.SOCKS5Proxy, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
		dialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	}

	client.httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: dialContext,
			Proxy:       http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				RootCAs:            roots,
				InsecureSkipVerify: client.InsecureSkipVerify,
			},
		},
	}

	return client, nil
}

func (client *Client) Trace(writer io.Writer, traceLevel TraceLevel) {
	client.traceWriter = writer
	client.traceLevel = traceLevel
}

func (client *Client) traceRequest(request *http.Request) error {
	if client.traceLevel > 0 {
		what, err := httputil.DumpRequest(request, client.traceLevel == TraceAll)
		if err != nil {
			return err
		}
		fmt.Fprintf(client.traceWriter, "---[ request ]------------------------------------------------------------------\n@C{%s}\n\n", what)
	}
	return nil
}

func (client *Client) traceResponse(response *http.Response) error {
	if client.traceLevel > 0 {
		what, err := httputil.DumpResponse(response, client.traceLevel == TraceAll)
		if err != nil {
			return err
		}
		fmt.Fprintf(client.traceWriter, "---[ response ]-----------------------------------------------------------------\n@W{%s}\n\n", what)
	}
	return nil
}

func (client *Client) url(path string) string {
	if path == "" || path[0:1] != "/" {
		path = "/" + path
	}
	scheme := client.Protocol
	if scheme == "" {
		scheme = "https"
	}

	if client.Bucket == "" {
		return fmt.Sprintf("%s://%s%s", scheme, client.Domain, path)
	}

	if client.UsePathBuckets {
		return fmt.Sprintf("%s://%s/%s%s", scheme, client.Domain, client.Bucket, path)
	} else {
		return fmt.Sprintf("%s://%s.%s%s", scheme, client.Bucket, client.Domain, path)
	}
}

func (client *Client) signature(req *http.Request, raw []byte) string {
	if client.SignatureVersion == 2 {
		return client.v2signature(req, raw)
	}
	if client.SignatureVersion == 4 {
		return client.v4signature(req, raw)
	}
	panic(fmt.Sprintf("unrecognized aws/s3 signature version %d", client.SignatureVersion))
}

func (client *Client) v2signature(req *http.Request, raw []byte) string {
	now := time.Now().UTC()

	req.Header.Set("x-amz-date", now.Format("20060102T150405Z"))
	req.Header.Set("host", regexp.MustCompile(`:.*`).ReplaceAllString(req.URL.Host, ""))
	if client.Token != "" {
		req.Header.Set("X-Amz-Security-Token", client.Token)
	}

	h := hmac.New(sha1.New, []byte(client.SecretAccessKey))
	h.Write([]byte(req.Method + "\n"))
	h.Write([]byte(req.Header.Get("Content-MD5") + "\n"))
	h.Write([]byte(req.Header.Get("Content-Type") + "\n"))
	h.Write([]byte(req.Header.Get("Date") + "\n"))
	h.Write(utils.V2Headers(req))
	h.Write(utils.V2Resource(client.Bucket, req))

	return fmt.Sprintf("AWS %s:%s", client.AccessKeyId, base64.StdEncoding.EncodeToString(h.Sum(nil)))
}

func (client *Client) v4signature(req *http.Request, raw []byte) string {
	/* step 0: assemble some temporary values we will need */
	now := time.Now().UTC()
	yyyymmdd := now.Format("20060102")
	scope := fmt.Sprintf("%s/%s/s3/aws4_request", yyyymmdd, client.Region)
	req.Header.Set("x-amz-date", now.Format("20060102T150405Z"))
	req.Header.Set("host", req.URL.Host)
	if client.Token != "" {
		req.Header.Set("X-Amz-Security-Token", client.Token)
	}

	payload := sha256.New()
	payload.Write(raw)
	hashed := hex.EncodeToString(payload.Sum(nil))
	req.Header.Set("x-amz-content-sha256", hashed)

	/* step 1: generate the CanonicalRequest (+sha256() it)

	   METHOD \n
	   uri() \n
	   querystring() \n
	   headers() \n
	   signed() \n
	   payload()
	*/

	headers, hexSignature := utils.V4Headers(req)
	canon := sha256.New()
	canon.Write([]byte(req.Method))
	canon.Write([]byte("\n"))
	canon.Write([]byte(url.PathEscape(req.URL.Path)))
	canon.Write([]byte("\n"))
	canon.Write(utils.V4QueryString(req.URL.RawQuery))
	canon.Write([]byte("\n"))
	canon.Write(hexSignature)
	canon.Write([]byte("\n"))
	canon.Write([]byte(hashed))

	//fmt.Printf("CANONICAL:\n---\n%s\n%s\n%s\n%s\n%s]---\n",
	//	req.Method, string(uriencode(req.URL.Path, false)), string(v4QueryString(req.URL.RawQuery)), string(hsig), hashed)

	/* step 2: generate the StringToSign

	   AWS4-HMAC-SHA256 \n
	   YYYYMMDDTHHMMSSZ \n
	   "yyyymmdd/region/s3/aws_request" \n
	   hex(sha256(canonical()))
	*/
	cleartext := "AWS4-HMAC-SHA256" +
		"\n" + now.Format("20060102T150405Z") +
		"\n" + scope +
		"\n" + hex.EncodeToString(canon.Sum(nil))

	//fmt.Printf("CLEARTEXT:\n---\n%s\n---\n", cleartext)

	/* step 3: generate the Signature

	   datekey = hmac-sha256("AWS4" + secret_key, YYYYMMDD)
	   datereg = hmac-sha256(datekey, region)
	   drsvc   = hmac-sha256(datereg, "s3")
	   sigkey  = hmac-sha256(drsvc, "aws4_request")

	   hex.EncodeToString(hmac-sha256(sigkey, cleartext))

	*/
	k1 := utils.Mac256([]byte("AWS4"+client.SecretAccessKey), []byte(yyyymmdd))
	k2 := utils.Mac256(k1, []byte(client.Region))
	k3 := utils.Mac256(k2, []byte("s3"))
	k4 := utils.Mac256(k3, []byte("aws4_request"))
	sig := hex.EncodeToString(utils.Mac256(k4, []byte(cleartext)))

	/* step 4: assemble and return the Authorize: header */
	return "AWS4-HMAC-SHA256" +
		" " + fmt.Sprintf("Credential=%s/%s", client.AccessKeyId, scope) +
		"," + fmt.Sprintf("SignedHeaders=%s", string(headers)) +
		"," + fmt.Sprintf("Signature=%s", sig)
}

func (client *Client) request(method, path string, payload []byte, headers *http.Header) (*http.Response, error) {
	in := bytes.NewBuffer(payload)
	req, err := http.NewRequest(method, client.url(path), in)
	if err != nil {
		return nil, err
	}

	/* copy in any headers */
	if headers != nil {
		for header, values := range *headers {
			for _, value := range values {
				req.Header.Add(header, value)
			}
		}
	}

	/* sign the request */
	req.ContentLength = int64(len(payload))
	req.Header.Set("Authorization", client.signature(req, payload))

	/* stupid continuation tokens sometimes have literal +'s in them */
	req.URL.RawQuery = regexp.MustCompile(`\+`).ReplaceAllString(req.URL.RawQuery, "%2B")

	/* optional debugging */
	if err := client.traceRequest(req); err != nil {
		return nil, err
	}

	/* submit the request */
	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	/* optional debugging */
	if err := client.traceResponse(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (client *Client) post(path string, payload []byte, headers *http.Header) (*http.Response, error) {
	return client.request("POST", path, payload, headers)
}

func (client *Client) put(path string, payload []byte, headers *http.Header) (*http.Response, error) {
	return client.request("PUT", path, payload, headers)
}

func (client *Client) get(path string, headers *http.Header) (*http.Response, error) {
	return client.request("GET", path, nil, headers)
}

func (client *Client) delete(path string, headers *http.Header) (*http.Response, error) {
	return client.request("DELETE", path, nil, headers)
}
