package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func readFileIntoString(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("Error read file: %w", err)
	}
	return string(data), err
}

func writeStringToFile(filename string, data string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error writeStringToFile: %w", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.WriteString(data)
	if err != nil {
		return fmt.Errorf("Error writeStringToFile: %w", err)
	}
	w.Flush()
	return err
}
