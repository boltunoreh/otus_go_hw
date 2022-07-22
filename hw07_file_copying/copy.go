package main

import (
	"errors"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, _ := os.Open(fromPath)
	toFile, _ := os.Create(toPath)

	fromFileInfo, _ := fromFile.Stat()
	fileSize := fromFileInfo.Size()

	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	if fileSize == 0 {
		return ErrUnsupportedFile
	}

	if limit == 0 {
		limit = fileSize
	} else if limit > (fileSize - offset) {
		limit = fileSize - offset
	}

	var barSize int64
	if limit > 0 {
		barSize = limit
	} else {
		barSize = fileSize - offset
	}

	bar := pb.Full.Start64(barSize)
	barReader := bar.NewProxyReader(fromFile)

	fromFile.Seek(offset, 0)
	io.CopyN(toFile, barReader, limit)

	bar.Finish()

	return nil
}
