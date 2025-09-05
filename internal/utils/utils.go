package utils

import (
	"io"
	"log"
)

// CloseOrLog closes the given io.Closer and logs an error if it occurs.
// The 'name' parameter is used to identify the resource being closed in the log message.
func CloseOrLog(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		log.Printf("Error while closing resource %q: %+v", name, err)
	}
}
