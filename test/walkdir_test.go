package test

import (
	"io/fs"
	"log"
	"os"
	"testing"
)

func TestWalkDir(t *testing.T) {
	root := "C:\\photo-patchouli\\temp"
	fileSystem := os.DirFS(root)
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		log.Println(path)
		return nil
	})
	if err != nil {
		return
	}
}
