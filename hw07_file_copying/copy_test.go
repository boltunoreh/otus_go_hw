package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tempFile, _ := os.CreateTemp("/tmp", "hw07_file_copying_*")

	err := Copy("testdata/input.txt", tempFile.Name(), 1000000000000000, 0)
	require.True(t, errors.Is(err, ErrOffsetExceedsFileSize))

	err = Copy("/dev/urandom", tempFile.Name(), 0, 0)
	require.True(t, errors.Is(err, ErrUnsupportedFile))
}
