package utils

import (
	"context"
	"fmt"
	"os"
	"time"
)

// FlagExists checks if file exists
func FlagExists(name string) bool {
	if len(name) == 0 {
		return false
	}
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// CreateFlag creates flag and returns true if no error
func CreateFlag(name string) bool {
	if len(name) == 0 {
		return false
	}
	var file, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

// WaitForFlag awaits the flag to appear in the file system
func WaitForFlag(name string, timeoutMS int, ctx context.Context) bool {

	checkEvery := time.NewTicker(300 * time.Millisecond)
	timeout := time.After(time.Duration(timeoutMS) * time.Millisecond)

	sp := NewSpinner(fmt.Sprintf(" waiting for the flag %s ..", name))
	defer sp.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-timeout:
			checkEvery.Stop()
			return false
		case <-checkEvery.C:
			if FlagExists(name) {
				return true
			}
		}
	}
}

// CreateFlag creates flag and returns true if no error
func CreateFlagWithContents(name, cont string) bool {
	if len(name) == 0 {
		return false
	}
	var file, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return false
	}
	file.WriteString(cont)

	defer file.Close()
	return true
}

// DeleteFlag returns false if there was no flag
func DeleteFlag(name string) bool {
	var err = os.Remove(name)
	return err == nil
}
