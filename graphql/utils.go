package main

import (
	"fmt"
	"log"
)

func handleError(context string, err error) error {
	if err == nil {
		return nil
	}
	log.Printf("[ERROR] %s: %v\n", context, err)
	return fmt.Errorf("%s: %w", context, err)
}
