package test

import (
	"photowatcher"
	"testing"
)

func TestDiff(t *testing.T) {
	oldSnap := photowatcher.NewFileSnapshot("C:\\photo-patchouli\\temp\\3")
	newSnap := photowatcher.NewFileSnapshot("C:\\photo-patchouli\\temp\\00\\新建文件夹")
	oldSnap.Diff(newSnap) // 2w张图片仅需0.03秒，高效，迅速
}
