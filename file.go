package main

import (
	"io"
	"log"
	"os"
)

// openFile opens a file based on the
// given path and returns a *os.File pointer to it
func openFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %q: %v\n", path, err)
		return nil, err
	}
	return file, nil
}

// closeFile takes in a *os.File pointer and closes the file
func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Printf("Error closing file %q: %v\n", file.Name(), err)
	}
	return
}

// readFile takes in a *os.File pointer, reads in the data,
// and then returns it as an array of bytes
func readFile(file *os.File) []byte {
	body, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file %q: %v\n", file.Name(), err)
	}
	return body
}

// openAndReadFile takes in a file path, opens the file,
// converts it to an array of bytes and returns the result
func openAndReadFile(path string) ([]byte, error) {
	file, err := openFile(path)
	if err != nil {
		log.Printf("Error opening file %q: %v\n", path, err)
		return nil, err
	}
	defer closeFile(file)
	content := readFile(file)
	return content, err
}
