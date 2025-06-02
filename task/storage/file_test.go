package storage

import (
	"testing"
)

func TestFileCreation(t *testing.T) {
	Append("Hello")
	Clean()

}
