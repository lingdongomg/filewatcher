package main

import "log"
import "photowatcher"

type MyCallback struct{}

func (cb *MyCallback) OnPathChanged(cbe photowatcher.CallBackEvent) {
	log.Println("CallBack:", cbe.Path, cbe.Op)
}

func main() {
	callback := &MyCallback{}
	photowatcher.SetPathCallback(callback)

	multiWatcher, _ := photowatcher.NewMultiWatcher()
	err := multiWatcher.Add("C:\\photo-patchouli\\temp")
	if err != nil {
		log.Println(err)
		return
	}
	select {}
	multiWatcher.Close()
}
