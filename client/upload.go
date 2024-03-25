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
	Key        string
	client     *Client
	partNumber int
	id         string
	signature  string
	path       string

	parts []xmlpart
}

func (upload *Upload) nextPart() int {
	upload.parts = append(upload.parts, xmlpart{})
	upload.partNumber = upload.partNumber + 1
	return upload.partNumber
}

func (upload *Upload) writePart(body []byte, partNumber int) error {
	if partNumber > 10000 {
		return fmt.Errorf("S3 limits the number of multipart upload segments to 10k")
	}

	res, err := upload.client.put(fmt.Sprintf("%s?partNumber=%d&uploadId=%s", upload.path, partNumber, upload.id), body, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	upload.parts[partNumber-1] = xmlpart{
		PartNumber: partNumber,
		ETag:       res.Header.Get("ETag"),
	}
	return nil
}

func (upload *Upload) Write(body []byte) error {
	return upload.writePart(body, upload.nextPart())
}

func (upload *Upload) Done() error {
	var payload struct {
		XMLName xml.Name  `xml:"CompleteMultipartUpload"`
		Parts   []xmlpart `xml:"Part"`
	}
	payload.Parts = upload.parts

	body, err := xml.Marshal(payload)
	if err != nil {
		return err
	}

	res, err := upload.client.post(fmt.Sprintf("%s?uploadId=%s", upload.path, upload.id), body, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return ResponseError(res)
	}
	return nil
}

func (upload *Upload) ParallelStream(in io.Reader, blockSize int, threads int) (int64, error) {
	if blockSize < 5*1024*1024 {
		return 0, fmt.Errorf("S3 requires block sizes of 5MB or higher")
	}

	type chunk struct {
		partNumber int
		block      []byte
	}

	var wg sync.WaitGroup
	chunks := make(chan chunk)
	errors := make(chan error, threads)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		//idx := i
		go func() {
			defer wg.Done()
			//fmt.Printf("io/%d: waiting for first chunk...\n", idx)
			for chunk := range chunks {
				//fmt.Printf("io/%d: writing %d bytes in chunk %d via writePart...\n", idx, len(chunk.block), chunk.n)
				err := upload.writePart(chunk.block, chunk.partNumber)
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
		buf := make([]byte, blockSize)
		readBytes, err := io.ReadAtLeast(in, buf, blockSize)
		if err != nil && err != io.ErrUnexpectedEOF {
			if err == io.EOF {
				return total, nil
			}
			return total, err
		}

		chunks <- chunk{
			partNumber: upload.nextPart(),
			block:      buf[0:readBytes],
		}

		total += int64(readBytes)
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

func (upload *Upload) Stream(reader io.Reader, blockSize int) (int64, error) {
	return upload.ParallelStream(reader, blockSize, 1)
}
