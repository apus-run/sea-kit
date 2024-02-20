package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetWorkDir return the current working directory
func GetWorkDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Cannot get the current directory: %v, using $HOME directory!", err)
		dir, err = os.UserHomeDir()
		if err != nil {
			log.Printf("Cannot get the user home directory: %v, using /tmp directory!", err)
			dir = os.TempDir()
		}
	}
	return dir
}

// MakeDirectory return the writeable filename
func MakeDirectory(filename string) string {
	dir, file := filepath.Split(filename)
	if len(dir) <= 0 {
		dir = GetWorkDir()
	}
	if len(file) <= 0 {
		return dir
	}
	if strings.HasPrefix(dir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Cannot get the user home directory: %v, using /tmp directory as home", err)
			home = os.TempDir()
		}
		dir = filepath.Join(home, dir[2:])
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Printf("Cannot get the absolute path: %v", err)
		dir = GetWorkDir()
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Printf("Cannot create the directory: %v", err)
			dir = GetWorkDir()
		}
	}

	return filepath.Join(dir, file)
}
