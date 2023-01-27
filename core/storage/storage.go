package storage

import "os"

// Client storage client
type Client interface {
	Upload(file *os.File, fileName string, contentType string) error
}
