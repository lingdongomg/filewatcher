package main

import "log"
import "photowatcher"

type MyCallback struct{}

func (cb *MyCallback) OnPathChanged(cbe photowatcher.CallBackEvent) {
	log.Println("[SDK的回调] ", cbe.Path, cbe.Op)
}

func main() {
	callback := &MyCallback{}
	photowatcher.SetPathCallback(callback)

	multiwatchar, _ := photowatcher.NewMultiWatcher()
	err := multiwatchar.Add("C:\\temp\\1")
	if err != nil {
		log.Println(err)
		return
	}
	select {}
	multiwatchar.Close()
}
