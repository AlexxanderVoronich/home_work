package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	isRemoveTempDir := true
	tests := []struct {
		sourcePath string
		destDir    string
		destFile   string
		checkFile  string
		errors     []error
		offset     int64
		limit      int64
	}{
		{
			sourcePath: "testdata/1.png",
			destDir:    "tmp",
			destFile:   "image.png",
			checkFile:  "testdata/1.png",
			offset:     0,
			limit:      0,
		},
		{
			sourcePath: "testdata/input.txt",
			destDir:    "tmp",
			destFile:   "1.txt",
			checkFile:  "testdata/out_offset0_limit0.txt",
			offset:     0,
			limit:      0,
		},
		{
			sourcePath: "testdata/input.txt",
			destDir:    "tmp",
			destFile:   "2.txt",
			checkFile:  "testdata/out_offset0_limit10.txt",
			offset:     0,
			limit:      10,
		},
		{
			sourcePath: "testdata/input.txt",
			destDir:    "tmp",
			destFile:   "3.txt",
			checkFile:  "testdata/out_offset0_limit1000.txt",
			offset:     0,
			limit:      1000,
		},
		{
			sourcePath: "testdata/input.txt",
			destDir:    "tmp",
			destFile:   "4.txt",
			checkFile:  "testdata/out_offset0_limit10000.txt",
			offset:     0,
			limit:      10000,
		},
		{
			sourcePath: "testdata/input.txt",
			destDir:    "tmp",
			destFile:   "5.txt",
			checkFile:  "testdata/out_offset100_limit1000.txt",
			offset:     100,
			limit:      1000,
		},
		{
			sourcePath: "testdata/input.txt",
			destDir:    "tmp",
			destFile:   "6.txt",
			checkFile:  "testdata/out_offset6000_limit1000.txt",
			offset:     6000,
			limit:      1000,
		},
		{
			sourcePath: "testdata/input.txt",
			destDir:    "tmp",
			destFile:   "7.txt",
			checkFile:  "Error",
			errors:     []error{ErrOffsetExceedsFileSize},
			offset:     10000,
			limit:      0,
		},
		{
			sourcePath: "/dev/urandom",
			destDir:    "tmp",
			destFile:   "8.txt",
			checkFile:  "Error",
			errors:     []error{ErrUnsupportedFile, ErrFileNotFound},
			offset:     0,
			limit:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.destFile, func(t *testing.T) {
			if tt.checkFile != "Error" {
				err := Copy(tt.sourcePath, filepath.Join(tt.destDir, tt.destFile), tt.offset, tt.limit)
				require.NoError(t, err, "Error copying file")

				checkContent, err := os.ReadFile(tt.checkFile)
				require.NoError(t, err, "Error reading check file")

				destContent, err := os.ReadFile(filepath.Join(tt.destDir, tt.destFile))
				require.NoError(t, err, "Error reading destination file")

				require.Equal(t, string(checkContent), string(destContent), "Content mismatch")
			} else {
				err := Copy(tt.sourcePath, filepath.Join(tt.destDir, tt.destFile), tt.offset, tt.limit)
				// require.EqualError(t, err, tt.err.Error())
				checkErr := false
				for _, e := range tt.errors {
					checkErr = checkErr || errors.Is(err, e)
				}
				require.True(t, checkErr)
			}
		})
	}

	if isRemoveTempDir {
		err := os.RemoveAll(tests[0].destDir)
		if err != nil {
			fmt.Printf("Error while deleting a folder: %v\n", err)
			return
		}
	}
}
