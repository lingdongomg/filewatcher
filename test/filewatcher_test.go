package test

import (
	"filewatcher"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

// 常量，测试监听路径，根据自己的文件系统更改
var testPath = `C:\photo-patchouli\temp`

type MyCallback struct{}

func (cb *MyCallback) OnPathChanged(cbe filewatcher.CallBackEvent) {
	log.Println("CallBack:", cbe.Path, cbe.Op, cbe.Size, cbe.ModTime)
}

func TestFileWatcher(t *testing.T) {
	callback := &MyCallback{}
	filewatcher.SetPathCallback(callback)
	multiWatcher, _ := filewatcher.NewMultiWatcher()
	err := multiWatcher.Add(testPath)
	if err != nil {
		log.Println(err)
		return
	}
	select {}
}

func TestClose(t *testing.T) {
	callback := &MyCallback{}
	filewatcher.SetPathCallback(callback)
	multiWatcher, _ := filewatcher.NewMultiWatcher()
	err := multiWatcher.Add(testPath)
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
	if _, err := os.Stat(testPath + "\\1"); os.IsNotExist(err) {
		// 如果文件不存在
		err := os.Mkdir(testPath+"\\1", os.ModePerm)
		if err != nil {
			// 如果创建文件夹失败
			panic(err)
		}
	}
	if _, err := os.Stat(testPath + "\\2"); os.IsNotExist(err) {
		// 如果文件不存在
		err := os.Mkdir(testPath+"\\2", os.ModePerm)
		if err != nil {
			// 如果创建文件夹失败
			panic(err)
		}
	}

	// 如果 C:\photo-patchouli\temp\1 存在文件1.jpg，则删除文件1.jpg
	// 如果 C:\photo-patchouli\temp\2 不存在文件2.jpg，则创建文件2.jpg
	file1 := testPath + "\\1\\1.jpg"
	if _, err := os.Stat(file1); err == nil {
		// 如果文件存在
		err := os.Remove(file1)
		if err != nil {
			// 如果删除文件失败
			panic(err)
		}
	}

	// 如果 C:\photo-patchouli\temp\2 不存在文件2.jpg，则创建文件2.jpg
	file2 := testPath + "\\2\\2.jpg"
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		// 如果文件不存在
		_, err := os.Create(file2)
		if err != nil {
			// 如果创建文件失败
			panic(err)
		}
	}

	callback := &MyCallback{}
	filewatcher.SetPathCallback(callback)
	multiWatcher, _ := filewatcher.NewMultiWatcher()
	err := multiWatcher.Add(testPath)
	if err != nil {
		log.Println(err)
		return
	}

	// 将文件1.jpg移动到C:\photo-patchouli\temp\2
	err = os.Rename(testPath+"\\2\\2.jpg", testPath+"\\1\\1.jpg")

	time.Sleep(7 * time.Second)
	err = multiWatcher.Close()
	if err != nil {
		log.Println(err)
		return
	}
}

func TestDirMove(t *testing.T) {
	// 如果C:\photo-patchouli\temp\3\ 和 C:\photo-patchouli\temp\4\ 不存在，则创建
	if _, err := os.Stat(testPath + "\\3"); os.IsNotExist(err) {
		// 如果文件不存在
		err := os.Mkdir(testPath+"\\3", os.ModePerm)
		if err != nil {
			// 如果创建文件夹失败
			panic(err)
		}
	}
	if _, err := os.Stat(testPath + "\\4"); os.IsNotExist(err) {
		// 如果文件不存在
		err := os.Mkdir(testPath+"\\4", os.ModePerm)
		if err != nil {
			// 如果创建文件夹失败
			panic(err)
		}
	}

	// 如果 C:\photo-patchouli\temp\3 存在文件夹4，则删除文件夹4
	if _, err := os.Stat(testPath + "\\3\\4"); err == nil {
		// 如果文件夹存在
		err := os.RemoveAll(testPath + "\\3\\4")
		if err != nil {
			// 如果删除文件夹失败
			panic(err)
		}

	}

	// 在文件夹4下创建10张jpg文件, 如果文件不存在，则创建
	for i := 0; i < 10; i++ {
		file := testPath + "\\4\\" + strconv.Itoa(i) + `.jpg`
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
	filewatcher.SetPathCallback(callback)
	multiWatcher, _ := filewatcher.NewMultiWatcher()
	err := multiWatcher.Add(testPath)
	if err != nil {
		log.Println(err)
		return
	}
	// 将文件夹4移动到C:\photo-patchouli\temp\3
	err = os.Rename(testPath+"\\4", testPath+"\\3\\4")

	time.Sleep(10 * time.Second)
	err = multiWatcher.Close()
	if err != nil {
		log.Println(err)
		return
	}
}

func TestAddPath(t *testing.T) {
	callback := &MyCallback{}
	filewatcher.SetPathCallback(callback)
	multiWatcher, _ := filewatcher.NewMultiWatcher()
	err := multiWatcher.Add(testPath + "\\1")
	if err != nil {
		log.Println(err)
		return
	}
	err = multiWatcher.Add(testPath + "\\1")
	if err != nil {
		return
	}

	err = multiWatcher.Add(testPath + "\\2")
	if err != nil {
		return
	}

	file1 := testPath + "\\1\\1.jpg"
	if _, err := os.Stat(file1); err == nil {
		// 如果文件存在
		err := os.Remove(file1)
		if err != nil {
			// 如果删除文件失败
			panic(err)
		}
	}

	// 如果 C:\photo-patchouli\temp\2 不存在文件2.jpg，则创建文件2.jpg
	file2 := testPath + "\\2\\2.jpg"
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		// 如果文件不存在
		file, err := os.Create(file2)
		err = file.Close()
		if err != nil {
			return
		}
		if err != nil {
			// 如果创建文件失败
			panic(err)
		}
	}

	// 将文件1.jpg移动到C:\photo-patchouli\temp\2
	err = os.Rename(testPath+"\\2\\2.jpg", testPath+"\\1\\1.jpg")
	if err != nil {
		log.Println(err)
	}

	time.Sleep(7 * time.Second)
	err = multiWatcher.Close()
	if err != nil {
		log.Println(err)
		return
	}
}

func TestClosePath(t *testing.T) {
	file1 := testPath + "\\1\\1.jpg"
	if _, err := os.Stat(file1); err == nil {
		// 如果文件存在
		err := os.Remove(file1)
		if err != nil {
			// 如果删除文件失败
			panic(err)
		}
	}

	// 如果 C:\photo-patchouli\temp\2 不存在文件2.jpg，则创建文件2.jpg
	file2 := testPath + "\\2\\2.jpg"
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		// 如果文件不存在
		file, err := os.Create(file2)
		err = file.Close()
		if err != nil {
			return
		}
		if err != nil {
			// 如果创建文件失败
			panic(err)
		}
	}

	callback := &MyCallback{}
	filewatcher.SetPathCallback(callback)
	multiWatcher, _ := filewatcher.NewMultiWatcher()
	err := multiWatcher.Add(testPath + "\\1")
	if err != nil {
		log.Println(err)
		return
	}

	err = multiWatcher.Add(testPath + "\\2")
	if err != nil {
		return
	}

	err = multiWatcher.Remove(testPath + "\\2")
	if err != nil {
		return
	}

	// 将文件1.jpg移动到C:\photo-patchouli\temp\2
	err = os.Rename(testPath+"\\2\\2.jpg", testPath+"\\1\\1.jpg")
	if err != nil {
		log.Println(err)
	}

	err = multiWatcher.Close()
	if err != nil {
		log.Println(err)
		return
	}

}
