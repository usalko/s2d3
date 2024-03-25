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
	"encoding/xml"
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

	"github.com/usalko/s2d3/models"
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

func (client *Client) v2signature(request *http.Request, raw []byte) string {
	now := time.Now().UTC()

	request.Header.Set("x-amz-date", now.Format("20060102T150405Z"))
	request.Header.Set("host", regexp.MustCompile(`:.*`).ReplaceAllString(request.URL.Host, ""))
	if client.Token != "" {
		request.Header.Set("X-Amz-Security-Token", client.Token)
	}

	hmacHash := hmac.New(sha1.New, []byte(client.SecretAccessKey))
	hmacHash.Write([]byte(request.Method + "\n"))
	hmacHash.Write([]byte(request.Header.Get("Content-MD5") + "\n"))
	hmacHash.Write([]byte(request.Header.Get("Content-Type") + "\n"))
	hmacHash.Write([]byte(request.Header.Get("Date") + "\n"))
	hmacHash.Write(utils.V2Headers(request))
	hmacHash.Write(utils.V2Resource(client.Bucket, request))

	return fmt.Sprintf("AWS %s:%s", client.AccessKeyId, base64.StdEncoding.EncodeToString(hmacHash.Sum(nil)))
}

func (client *Client) v4signature(request *http.Request, raw []byte) string {
	/* step 0: assemble some temporary values we will need */
	now := time.Now().UTC()
	yyyymmdd := now.Format("20060102")
	scope := fmt.Sprintf("%s/%s/s3/aws4_request", yyyymmdd, client.Region)
	request.Header.Set("x-amz-date", now.Format("20060102T150405Z"))
	request.Header.Set("host", request.URL.Host)
	if client.Token != "" {
		request.Header.Set("X-Amz-Security-Token", client.Token)
	}

	payload := sha256.New()
	payload.Write(raw)
	hashed := hex.EncodeToString(payload.Sum(nil))
	request.Header.Set("x-amz-content-sha256", hashed)

	/* step 1: generate the CanonicalRequest (+sha256() it)

	   METHOD \n
	   uri() \n
	   querystring() \n
	   headers() \n
	   signed() \n
	   payload()
	*/

	headers, hexSignature := utils.V4Headers(request)
	canon := sha256.New()
	canon.Write([]byte(request.Method))
	canon.Write([]byte("\n"))
	canon.Write([]byte(url.PathEscape(request.URL.Path)))
	canon.Write([]byte("\n"))
	canon.Write(utils.V4QueryString(request.URL.RawQuery))
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

// --------------------------------------------------------------------------------------------

func (client *Client) List() ([]models.Object, error) {
	objects := make([]models.Object, 0)
	clientToken := ""
	for {
		response, err := client.get(fmt.Sprintf("/?list-type=2&fetch-owner=true%s", clientToken), nil)
		if err != nil {
			return nil, err
		}

		var parsedBody struct {
			XMLName  xml.Name `xml:"ListBucketResult"`
			Next     string   `xml:"NextContinuationToken"`
			Contents []struct {
				Key          string `xml:"Key"`
				LastModified string `xml:"LastModified"`
				ETag         string `xml:"ETag"`
				Size         int64  `xml:"Size"`
				StorageClass string `xml:"StorageClass"`
				Owner        struct {
					ID          string `xml:"ID"`
					DisplayName string `xml:"DisplayName"`
				} `xml:"Owner"`
			} `xml:"Contents"`
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		if response.StatusCode != 200 {
			return nil, ResponseErrorFrom(body)
		}

		err = xml.Unmarshal(body, &parsedBody)
		if err != nil {
			return nil, err
		}

		for _, f := range parsedBody.Contents {
			mod, _ := time.Parse("2006-01-02T15:04:05.000Z", f.LastModified)
			objects = append(objects, models.Object{
				Key:          f.Key,
				LastModified: mod,
				ETag:         f.ETag[1 : len(f.ETag)-1],
				Size:         utils.SizeInBytes(f.Size),
				StorageClass: f.StorageClass,
				OwnerID:      f.Owner.ID,
				OwnerName:    f.Owner.DisplayName,
			})
		}

		if parsedBody.Next == "" {
			return objects, nil
		}

		clientToken = fmt.Sprintf("&continuation-token=%s", parsedBody.Next)
	}
}

// --------------------------------------------------------------------------------------------

func (client *Client) GetACL(key string) ([]models.Grant, error) {
	res, err := client.get(key+"?acl", nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, ResponseError(res)
	}

	var r struct {
		XMLName xml.Name `xml:"AccessControlPolicy"`
		List    struct {
			Grant []struct {
				Grantee struct {
					ID   string `xml:"ID"`
					Name string `xml:"DisplayName"`
					URI  string `xml:"URI"`
				} `xml:"Grantee"`
				Permission string `xml:"Permission"`
			} `xml:"Grant"`
		} `xml:"AccessControlList"`
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := xml.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	var acl []models.Grant
	for _, g := range r.List.Grant {
		group := ""
		if g.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" {
			group = "EVERYONE"
		}
		acl = append(acl, models.Grant{
			GranteeID:   g.Grantee.ID,
			GranteeName: g.Grantee.Name,
			Group:       group,
			Permission:  g.Permission,
		})
	}
	return acl, nil
}

func (client *Client) ChangeACL(path, acl string) error {
	headers := make(http.Header)
	headers.Set("x-amz-acl", acl)

	res, err := client.put(path+"?acl", nil, &headers)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return ResponseError(res)
	}

	return nil
}

// --------------------------------------------------------------------------------------------

func (client *Client) CreateBucket(name, region, acl string) error {
	/* validate that the bucket name is:

	   - between 3 and 63 characters long (inclusive)
	   - not include periods (for TLS wildcard matching)
	   - lower case
	   - rfc952 compliant
	*/
	if ok, _ := regexp.MatchString(`^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$`, name); !ok {
		return fmt.Errorf("invalid s3 bucket name")
	}

	was := client.Bucket
	defer func() { client.Bucket = was }()
	client.Bucket = name

	body := []byte{}
	if region != "" {
		var payload struct {
			XMLName xml.Name `xml:"CreateBucketConfiguration"`
			Region  string   `xml:"LocationConstraint"`
		}
		payload.Region = region

		var err error
		body, err = xml.Marshal(payload)
		if err != nil {
			return err
		}
	}

	headers := make(http.Header)
	headers.Set("x-amz-acl", acl)

	res, err := client.put("/", body, &headers)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return ResponseError(res)
	}

	return nil
}

func (client *Client) DeleteBucket(name string) error {
	was := client.Bucket
	defer func() { client.Bucket = was }()
	client.Bucket = name

	response, err := client.delete("/", nil)
	if err != nil {
		return err
	}

	if response.StatusCode != 204 {
		return ResponseError(response)
	}

	return nil
}

func (client *Client) ListBuckets() ([]models.Bucket, error) {
	prev := client.Bucket
	client.Bucket = ""
	response, err := client.get("/", nil)
	client.Bucket = prev
	if err != nil {
		return nil, err
	}

	var request struct {
		XMLName xml.Name `xml:"ListAllMyBucketsResult"`
		Owner   struct {
			ID          string `xml:"ID"`
			DisplayName string `xml:"DisplayName"`
		} `xml:"Owner"`
		Buckets struct {
			Bucket []struct {
				Name         string `xml:"Name"`
				CreationDate string `xml:"CreationDate"`
			} `xml:"Bucket"`
		} `xml:"Buckets"`
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ResponseErrorFrom(body)
	}

	err = xml.Unmarshal(body, &request)
	if err != nil {
		return nil, err
	}

	result := make([]models.Bucket, len(request.Buckets.Bucket))
	for i, bkt := range request.Buckets.Bucket {
		result[i].OwnerID = request.Owner.ID
		result[i].OwnerName = request.Owner.DisplayName
		result[i].Name = bkt.Name

		created, _ := time.Parse("2006-01-02T15:04:05.000Z", bkt.CreationDate)
		result[i].CreationDate = created
	}
	return result, nil
}

// --------------------------------------------------------------------------------------------

func (client *Client) Delete(path string) error {
	res, err := client.delete(path, nil)
	if err != nil {
		return err
	}

	if res.StatusCode != 204 {
		return ResponseError(res)
	}

	return nil
}

// --------------------------------------------------------------------------------------------

func (client *Client) Get(key string) (io.Reader, error) {
	res, err := client.get(key, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, ResponseError(res)
	}

	return res.Body, nil
}

// --------------------------------------------------------------------------------------------

func (client *Client) NewUpload(path string, headers *http.Header) (*Upload, error) {
	res, err := client.post(path+"?uploads", nil, headers)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, ResponseErrorFrom(b)
	}

	var payload struct {
		Bucket   string `xml:"Bucket"`
		Key      string `xml:"Key"`
		UploadId string `xml:"UploadId"`
	}
	err = xml.Unmarshal(b, &payload)
	if err != nil {
		return nil, err
	}

	return &Upload{
		Key: payload.Key,

		client:     client,
		id:         payload.UploadId,
		path:       path,
		partNumber: 0,
	}, nil
}
