package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		err = dst.Close()
	}()

	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		err = src.Close()
	}()

	srcInfo, err := src.Stat()
	if err != nil {
		return err
	}
	if !srcInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := srcInfo.Size()

	if offset > 0 {
		if offset > fileSize {
			return ErrOffsetExceedsFileSize
		}
		_, err = src.Seek(offset, 0)
		if err != nil {
			return err
		}
		fileSize -= offset
	}

	var reader io.Reader = src
	if limit > 0 {
		reader = io.LimitReader(src, limit)
		if fileSize > limit {
			fileSize = limit
		}
	}

	bar := pb.Full.Start64(fileSize)
	barReader := bar.NewProxyReader(reader)

	_, err = io.Copy(dst, barReader)

	bar.Finish()
	return err
}
