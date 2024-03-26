package services

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/usalko/s2d3/utils"
)

var BREAKPOINTS = [...]utils.SizeInBytes{
	4096,
	16384,
	65536,
	4194304,
	16777216,
	4294967296,
}
var BREAKPOINTS_DELTA = [...]utils.SizeInBytes{
	0,
	4096,
	16384,
	65536,
	4194304,
	16777216,
}

type Storage struct {
	RootFolder string
}

func findBreakpoint(dataSize int) (int, int) {
	breakpointIndex := 0
	countOfSegments := 1
	for index, breakpoint := range BREAKPOINTS {
		if int(breakpoint) > dataSize {
			if index > 0 {
				breakpointIndex = index - 1
				countOfSegments = dataSize / int(BREAKPOINTS[breakpointIndex])
			}
			break
		}
		if int(breakpoint+BREAKPOINTS_DELTA[index]) > dataSize {
			breakpointIndex = index
			countOfSegments = 1
			break
		}
	}
	return breakpointIndex, countOfSegments
}

func (storage *Storage) Init() {
	os.Mkdir(storage.RootFolder, fs.ModeDir|0775)
}

func (storage *Storage) CheckUpload(bucketName string, objectKey string, suffix string, uploadDone UploadDone) error {
	file, err := os.Open(strings.Join([]string{
		storage.RootFolder,
		bucketName,
		objectKey,
	}, "/"))
	defer file.Close()

	if err != nil {
		return err
	}
	return nil
}

func (storage *Storage) PushData(bucketName string, objectKey string, suffix string, reader io.ReadCloser) error {
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	breakpointIndex, countOfSegments := findBreakpoint(len(content))
	segmentSize := BREAKPOINTS[breakpointIndex]
	segmentDelta := BREAKPOINTS_DELTA[breakpointIndex]

	objectParentPath := strings.Join([]string{
		storage.RootFolder,
		bucketName,
	}, "/")
	err = os.MkdirAll(objectParentPath, fs.ModeDir|0775)
	if err != nil {
		return err
	}

	if countOfSegments == 1 {
		// Save object as single file

		err = os.WriteFile(strings.Join([]string{
			objectParentPath,
			objectKey,
		}, "/"), content, 0644)
		if err != nil {
			return err
		}

	} else if countOfSegments > 1 {
		println(segmentSize)
		println(segmentDelta)
		return fmt.Errorf("not implemented for countOfSegments %d", countOfSegments)
	} else {
		return fmt.Errorf("invalid count of segments for content length %d", len(content))
	}
	return nil
}

func (storage *Storage) GetData(bucketName string, objectKey string, suffix string) ([]byte, error) {
	file, err := os.Open(strings.Join([]string{
		storage.RootFolder,
		bucketName,
		objectKey,
	}, "/"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("can't open object %s/%s cause not implemented", bucketName, objectKey)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
