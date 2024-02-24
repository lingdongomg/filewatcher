package photowatcher

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"time"
)

type MultiWatcher struct {
	watcher   *fsnotify.Watcher
	paths     map[string]bool
	snapshots map[string]*FileSnapshot
}

func NewMultiWatcher() (*MultiWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	multichar := &MultiWatcher{
		watcher:   watcher,
		paths:     make(map[string]bool),
		snapshots: make(map[string]*FileSnapshot),
	}
	go fileChangeHandler(multichar)
	return multichar, nil
}

func fileChangeHandler(multichar *MultiWatcher) {
	pending := false
	eventCh := make(chan bool, 1)
	var timer *time.Timer
	for {
		select {
		case _, ok := <-multichar.watcher.Events:
			if !ok {
				log.Println("关闭监听事件管道")
				return
			}

			if !pending {
				pending = true
				timer = time.AfterFunc(5*time.Second, func() {
					eventCh <- true
				})
			} else {
				timer.Stop()
				timer = time.AfterFunc(5*time.Second, func() {
					eventCh <- true
				})
			}

		case err, ok := <-multichar.watcher.Errors:
			if !ok {
				log.Println("关闭监听错误管道")
				return
			}
			log.Println("监听错误管道发现:", err)

		case <-eventCh:
			pending = false
			log.Println("文件变化事件处理")
			for path := range multichar.paths {
				snapshot := multichar.snapshots[path]
				newSnapshot := NewFileSnapshot(path)
				diffs := snapshot.Diff(newSnapshot)
				if len(diffs) > 0 {
					multichar.snapshots[path] = newSnapshot
					for _, diff := range diffs {
						log.Println("文件变化：", diff)
						if cb != nil {
							cbe := CallBackEvent{Path: diff.AbsPath, Op: diff.Op}
							cb.OnPathChanged(cbe)
						} else {
							log.Println("没有回调函数")
						}
						// 如果diff.AbsPath是文件夹，将其添加到watcher的监听路径
						fileInfo, err := os.Stat(diff.AbsPath)
						if err != nil {
							log.Println("获取文件信息失败:", err)
							continue
						}
						if fileInfo.IsDir() {
							err := multichar.watcher.Add(diff.AbsPath)
							if err != nil {
								log.Println("添加目录到监听路径失败:", err)
							}
						}
					}
				}
			}
		}
	}
}

func (mw *MultiWatcher) Add(path string) error {
	if _, ok := mw.paths[path]; ok {
		return nil
	}
	err := mw.watcher.Add(path)
	if err != nil {
		return err
	}
	mw.paths[path] = true
	mw.snapshots[path] = NewFileSnapshot(path)
	traverseDir(mw.watcher, path)
	return nil
}

func traverseDir(watcher *fsnotify.Watcher, path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		fp := filepath.Join(path, file.Name())
		if file.IsDir() {
			err := watcher.Add(fp)
			log.Println("添加目录：", fp)
			if err != nil {
				return
			}
			traverseDir(watcher, fp)
		}
	}
}

func (mw *MultiWatcher) Remove(path string) error {
	if _, ok := mw.paths[path]; !ok {
		return nil
	}
	err := mw.watcher.Remove(path)
	if err != nil {
		return err
	}
	delete(mw.paths, path)
	delete(mw.snapshots, path)
	return nil
}

func (mw *MultiWatcher) Close() error {
	for path := range mw.paths {
		err := mw.Remove(path)
		if err != nil {
			log.Println("error removing", path, ":", err)
		}
	}
	return mw.watcher.Close()
}
