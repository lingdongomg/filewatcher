package test

import (
	"filewatcher"
	"log"
	"testing"
)

func TestDiff(t *testing.T) {
	oldSnap := filewatcher.NewFileSnapshot("C:\\photo-patchouli\\temp\\3\\1")
	newSnap := filewatcher.NewFileSnapshot("C:\\photo-patchouli\\temp\\3\\photos")
	// oldSnap.Diff(newSnap) // 2w张图片仅需0.03秒，高效，迅速
	for _, diff := range oldSnap.Diff(newSnap) {
		log.Println("Diff:", diff)
	}
}
