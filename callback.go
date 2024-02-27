package filewatcher

type CallBackEvent struct {
	Path    string
	Op      int
	Size    int64
	ModTime int64
}

type IPathCallback interface {
	OnPathChanged(cbe CallBackEvent)
}

var cb IPathCallback

func SetPathCallback(_cb IPathCallback) {
	cb = _cb
}
