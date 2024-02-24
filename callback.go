package photowatcher

type CallBackEvent struct {
	Path string
	Op   int
}

type IPathCallback interface {
	OnPathChanged(cbe CallBackEvent)
}

var cb IPathCallback

func SetPathCallback(_cb IPathCallback) {
	cb = _cb
}
