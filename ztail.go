package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "Usage: go run . -c <byte_count> <file1> [file2 ...]")
		os.Exit(1)
	}

	if os.Args[1] != "-c" {
		fmt.Fprintln(os.Stderr, "Only the -c option is supported.")
		os.Exit(1)
	}

	byteCount, err := strconv.Atoi(os.Args[2])
	if err != nil || byteCount <= 0 {
		fmt.Fprintln(os.Stderr, "Invalid byte count. Must be a positive integer.")
		os.Exit(1)
	}

	files := os.Args[3:]
	exitStatus := 0

	for i, filePath := range files {
		if i > 0 {
			fmt.Println()
		}

		fmt.Printf("==> %s <==\n", filePath)
		content, err := tailFile(filePath, byteCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filePath, err)
			exitStatus = 1
			continue
		}
		fmt.Print(content)
	}

	os.Exit(exitStatus)
}

func tailFile(filePath string, byteCount int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()

	offset := fileSize - int64(byteCount)
	if offset < 0 {
		offset = 0
	}

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, byteCount)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	return string(buffer[:n]), nil
}
