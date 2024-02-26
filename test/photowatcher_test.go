package test

import (
	"log"
	"os"
	"photowatcher"
	"strconv"
	"testing"
	"time"
)

type MyCallback struct{}

func (cb *MyCallback) OnPathChanged(cbe photowatcher.CallBackEvent) {
	log.Println("CallBack:", cbe.Path, cbe.Op)
}

func TestPhotoWatcher(t *testing.T) {
	callback := &MyCallback{}
	photowatcher.SetPathCallback(callback)
	multiWatcher, _ := photowatcher.NewMultiWatcher()
	err := multiWatcher.Add("C:\\photo-patchouli\\temp")
	if err != nil {
		log.Println(err)
		return
	}
	select {}
}

func TestClose(t *testing.T) {
	callback := &MyCallback{}
	photowatcher.SetPathCallback(callback)
	multiWatcher, _ := photowatcher.NewMultiWatcher()
	err := multiWatcher.Add("C:\\photo-patchouli\\temp")
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(5 * time.Second)
	err = multiWatcher.Close()
	if err != nil {
		log.Println(err)
		return
	}
}

func TestFileMove(t *testing.T) {
	// 如果C:\photo-patchouli\temp\1\ 和 C:\photo-patchouli\temp\2\ 不存在，则创建
	if _, err := os.Stat(`C:\photo-patchouli\temp\1`); os.IsNotExist(err) {
		// 如果文件不存在
		err := os.Mkdir(`C:\photo-patchouli\temp\1`, os.ModePerm)
		if err != nil {
			// 如果创建文件夹失败
			panic(err)
		}
	}
	if _, err := os.Stat(`C:\photo-patchouli\temp\2`); os.IsNotExist(err) {
		// 如果文件不存在
		err := os.Mkdir(`C:\photo-patchouli\temp\2`, os.ModePerm)
		if err != nil {
			// 如果创建文件夹失败
			panic(err)
		}
	}

	// 如果 C:\photo-patchouli\temp\1 存在文件1.jpg，则删除文件1.jpg
	// 如果 C:\photo-patchouli\temp\2 不存在文件2.jpg，则创建文件2.jpg
	file1 := `C:\photo-patchouli\temp\1\1.jpg`
	if _, err := os.Stat(file1); err == nil {
		// 如果文件存在
		err := os.Remove(file1)
		if err != nil {
			// 如果删除文件失败
			panic(err)
		}
	}

	// 如果 C:\photo-patchouli\temp\2 不存在文件2.jpg，则创建文件2.jpg
	file2 := `C:\photo-patchouli\temp\2\2.jpg`
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		// 如果文件不存在
		_, err := os.Create(file2)
		if err != nil {
			// 如果创建文件失败
			panic(err)
		}
	}

	callback := &MyCallback{}
	photowatcher.SetPathCallback(callback)
	multiWatcher, _ := photowatcher.NewMultiWatcher()
	err := multiWatcher.Add("C:\\photo-patchouli\\temp")
	if err != nil {
		log.Println(err)
		return
	}

	// 将文件1.jpg移动到C:\photo-patchouli\temp\2
	err = os.Rename(`C:\photo-patchouli\temp\2\2.jpg`, `C:\photo-patchouli\temp\1\1.jpg`)

	time.Sleep(6 * time.Second)
	err = multiWatcher.Close()
	if err != nil {
		log.Println(err)
		return
	}
}

func TestDirMove(t *testing.T) {
	// 如果C:\photo-patchouli\temp\3\ 和 C:\photo-patchouli\temp\4\ 不存在，则创建
	if _, err := os.Stat(`C:\photo-patchouli\temp\3`); os.IsNotExist(err) {
		// 如果文件不存在
		err := os.Mkdir(`C:\photo-patchouli\temp\3`, os.ModePerm)
		if err != nil {
			// 如果创建文件夹失败
			panic(err)
		}
	}
	if _, err := os.Stat(`C:\photo-patchouli\temp\4`); os.IsNotExist(err) {
		// 如果文件不存在
		err := os.Mkdir(`C:\photo-patchouli\temp\4`, os.ModePerm)
		if err != nil {
			// 如果创建文件夹失败
			panic(err)
		}
	}

	// 如果 C:\photo-patchouli\temp\3 存在文件夹4，则删除文件夹4
	if _, err := os.Stat(`C:\photo-patchouli\temp\3\4`); err == nil {
		// 如果文件夹存在
		err := os.RemoveAll(`C:\photo-patchouli\temp\3\4`)
		if err != nil {
			// 如果删除文件夹失败
			panic(err)
		}

	}

	// 在文件夹4下创建10张jpg文件, 如果文件不存在，则创建
	for i := 0; i < 10; i++ {
		file := `C:\photo-patchouli\temp\4\` + strconv.Itoa(i) + `.jpg`
		// 如果file不存在
		if _, err := os.Stat(file); os.IsNotExist(err) {
			_, err := os.Create(file)
			if err != nil {
				// 如果创建文件失败
				panic(err)
			}
		}

	}
	callback := &MyCallback{}
	photowatcher.SetPathCallback(callback)
	multiWatcher, _ := photowatcher.NewMultiWatcher()
	err := multiWatcher.Add("C:\\photo-patchouli\\temp")
	if err != nil {
		log.Println(err)
		return
	}
	// 将文件夹4移动到C:\photo-patchouli\temp\3
	err = os.Rename(`C:\photo-patchouli\temp\4\`, `C:\photo-patchouli\temp\3\4`)

	time.Sleep(10 * time.Second)
	err = multiWatcher.Close()
	if err != nil {
		log.Println(err)
		return
	}
}
