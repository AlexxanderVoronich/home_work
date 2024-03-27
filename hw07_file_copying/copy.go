package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb"
)

var (
	// ErrUnsupportedFile is returned when attempting to process an unsupported file type.
	ErrUnsupportedFile = errors.New("unsupported file")
	// ErrOffsetExceedsFileSize is returned when the given offset exceeds the file size.
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	// ErrFileNotFound is returned when a file is not found.
	ErrFileNotFound = errors.New("file not found")
	// ErrForTest is an error used for testing purposes only.
	ErrForTest = errors.New("test")
)

func validateFilePath(filePath string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("%w: %s", ErrFileNotFound, filePath)
	}
	return filePath, nil
}

// Copy copies the content from the source path to the destination path,
// with the specified offset and limit.
func Copy(fromPath, toPath string, offset, limit int64) error {
	// Check the source file and open it for reading
	fromPath, err := validateFilePath(fromPath)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	// Create the destination directory if it doesn't exist
	destDir := filepath.Dir(toPath)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err := os.MkdirAll(destDir, 0o755)
		if err != nil {
			return err
		}
	}

	// Create or truncate the destination file
	toFile, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer toFile.Close()

	// Set file pointer to offset
	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// Initialize progress variables
	copiedBytes := int64(0)
	bufferSize := int64(1024)
	isEndByLimit := false
	// create and start new bar
	all := fileInfo.Size()
	if limit > 0 {
		all = limit
	}
	bar := pb.StartNew(int(all))
	defer func() {
		bar.Finish()
	}()

	// Copy data in chunks
	buffer := make([]byte, bufferSize)
	for copiedBytes <= limit || limit == 0 {
		n, err := fromFile.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		// Adjust the chunk size if remaining bytes are less than the buffer size
		if limit > 0 && copiedBytes+int64(n) > limit {
			n = int(limit - copiedBytes)
			isEndByLimit = true
		}

		_, err = toFile.Write(buffer[:n])
		if err != nil {
			return err
		}
		copiedBytes += int64(n)

		// Calculate progress percentage
		// progress := float64(copiedBytes) / float64(limit) * 100
		// fmt.Printf("\rCopying %.2f%%", progress)
		bar.Set(int(copiedBytes))

		if isEndByLimit {
			break
		}
	}

	fmt.Println("\nFile copied successfully!")
	return nil
}
