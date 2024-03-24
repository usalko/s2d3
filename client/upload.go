package client

import (
	"encoding/xml"
	"fmt"
	"io"
	"sync"
)

type xmlpart struct {
	PartNumber int    `xml:"PartNumber"`
	ETag       string `xml:"ETag"`
}

type Upload struct {
	Key  string
	c    *Client
	n    int
	id   string
	sig  string
	path string

	parts []xmlpart
}

func (u *Upload) nextPart() int {
	u.parts = append(u.parts, xmlpart{})
	u.n = u.n + 1
	return u.n
}

func (upload *Upload) writePart(b []byte, n int) error {
	if n > 10000 {
		return fmt.Errorf("S3 limits the number of multipart upload segments to 10k")
	}

	res, err := upload.c.put(fmt.Sprintf("%s?partNumber=%d&uploadId=%s", upload.path, n, upload.id), b, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	upload.parts[n-1] = xmlpart{
		PartNumber: n,
		ETag:       res.Header.Get("ETag"),
	}
	return nil
}

func (upload *Upload) Write(b []byte) error {
	return upload.writePart(b, upload.nextPart())
}

func (upload *Upload) Done() error {
	var payload struct {
		XMLName xml.Name  `xml:"CompleteMultipartUpload"`
		Parts   []xmlpart `xml:"Part"`
	}
	payload.Parts = upload.parts

	b, err := xml.Marshal(payload)
	if err != nil {
		return err
	}

	res, err := upload.c.post(fmt.Sprintf("%s?uploadId=%s", upload.path, upload.id), b, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return ResponseError(res)
	}
	return nil
}

func (upload *Upload) ParallelStream(in io.Reader, block int, threads int) (int64, error) {
	if block < 5*1024*1024 {
		return 0, fmt.Errorf("S3 requires block sizes of 5MB or higher")
	}

	type chunk struct {
		n     int
		block []byte
	}

	var wg sync.WaitGroup
	chunks := make(chan chunk, 0)
	errors := make(chan error, threads)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		//idx := i
		go func() {
			defer wg.Done()
			//fmt.Printf("io/%d: waiting for first chunk...\n", idx)
			for chunk := range chunks {
				//fmt.Printf("io/%d: writing %d bytes in chunk %d via writePart...\n", idx, len(chunk.block), chunk.n)
				err := upload.writePart(chunk.block, chunk.n)
				if err != nil {
					errors <- err
					return
				}
				//fmt.Printf("io/%d: waiting for next chunk...\n", idx)
			}
		}()
	}

	var total int64
	defer func() {
		//fmt.Printf("... closing chunks channel\n")
		close(chunks)
		//fmt.Printf("... waiting for upload goroutines to exit")
		wg.Wait()
	}()

	for {
		buf := make([]byte, block)
		nread, err := io.ReadAtLeast(in, buf, block)
		if err != nil && err != io.ErrUnexpectedEOF {
			if err == io.EOF {
				return total, nil
			}
			return total, err
		}

		chunks <- chunk{
			n:     upload.nextPart(),
			block: buf[0:nread],
		}

		total += int64(nread)
		if err == io.ErrUnexpectedEOF {
			return total, nil
		}

		select {
		default:
		case err := <-errors:
			return total, err
		}
	}
}

func (upload *Upload) Stream(reader io.Reader, block int) (int64, error) {
	return upload.ParallelStream(reader, block, 1)
}
