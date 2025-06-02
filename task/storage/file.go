package storage

import (
	"log"
	"os"
)

const (
	outputFile = "timer_log.txt"
)

func Append(s string) error {

	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if _, err := f.Write([]byte(s + "\n")); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func Clean() error {
	if err := os.Remove(outputFile); err != nil {
		return err
	}
	return nil
}
