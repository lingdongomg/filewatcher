package main

import (
	"filewatcher"
	"log"
)

type MyCallback struct{}

func (cb *MyCallback) OnPathChanged(cbe filewatcher.CallBackEvent) {
	log.Println("CallBack:", cbe.Path, cbe.Op, cbe.Size, cbe.ModTime)
}

func main() {
	callback := &MyCallback{}
	filewatcher.SetPathCallback(callback)

	multiWatcher, _ := filewatcher.NewMultiWatcher()
	err := multiWatcher.Add("C:\\photo-patchouli\\temp")
	if err != nil {
		log.Println(err)
		return
	}
	select {}
}
