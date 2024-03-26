package services

import (
	"io"

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
			}
			countOfSegments = dataSize / int(BREAKPOINTS[breakpointIndex])
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

func (storage *Storage) PushData(bucketName string, objectKey string, suffix string, reader io.ReadCloser) error {
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	breakpointIndex, countOfSegments := findBreakpoint(len(content))

	println(breakpointIndex)
	println((countOfSegments))

	return nil
}
