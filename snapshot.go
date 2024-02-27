package filewatcher

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileNode struct {
	Name     string
	AbsPath  string
	Children []*FileNode
	IsDir    bool
	Size     int64
	ModTime  int64
}

type FileSnapshot struct {
	Root *FileNode
}

func newFileNode(name string, absPath string, isDir bool, size int64, modtime int64) *FileNode {
	return &FileNode{Name: name, AbsPath: absPath, IsDir: isDir, Size: size, ModTime: modtime}
}

func NewFileSnapshot(rootPath string) *FileSnapshot {
	root := newFileNode(string(os.PathSeparator), rootPath, true, 0, 0)
	buildTree(root, rootPath)
	return &FileSnapshot{Root: root}
}

func buildTree(node *FileNode, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Println(err)
		return
	}

	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			log.Println(err)
			continue
		}
		childNode := newFileNode(entry.Name(), childPath, entry.IsDir(), info.Size(), info.ModTime().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)))
		node.Children = append(node.Children, childNode)

		if entry.IsDir() {
			buildTree(childNode, childPath)
		}
	}
}

type Diff struct {
	AbsPath string
	Op      int // 0: file/directory deleted, 1: new file/directory, 2: file modified
	Size    int64
	ModTime int64
}

func (fs *FileSnapshot) Diff(fs2 *FileSnapshot) []Diff {
	return diffNodes(fs.Root, fs2.Root)
}

func diffNodes(oldNode *FileNode, newNode *FileNode) []Diff {
	if oldNode == nil && newNode == nil {
		return []Diff{}
	}
	if oldNode == nil {
		return handleNode(newNode, 1)
	}
	if newNode == nil {
		return handleNode(oldNode, 0)
	}

	oldChildren := make(map[string]*FileNode, len(oldNode.Children))
	for _, child := range oldNode.Children {
		oldChildren[child.Name] = child
	}

	newChildren := make(map[string]*FileNode, len(newNode.Children))
	for _, child := range newNode.Children {
		newChildren[child.Name] = child
	}

	return handleChildren(oldChildren, newChildren)
}

func handleNode(node *FileNode, op int) []Diff {
	var diffs []Diff
	absPath := node.AbsPath
	for _, child := range node.Children {
		var childDiffs []Diff
		if op == 0 {
			childDiffs = diffNodes(child, nil)
		} else if op == 1 {
			childDiffs = diffNodes(nil, child)
		}
		diffs = append(diffs, childDiffs...)
	}
	diffs = append(diffs, Diff{
		AbsPath: absPath,
		Op:      op,
		Size:    node.Size,
		ModTime: node.ModTime,
	})
	return diffs
}

func handleChildren(oldChildren map[string]*FileNode, newChildren map[string]*FileNode) []Diff {
	var diffs []Diff
	for name, oldChild := range oldChildren {
		newChild, ok := newChildren[name]
		if ok {
			if !oldChild.IsDir && !newChild.IsDir {
				if oldChild.Size != newChild.Size || oldChild.ModTime != newChild.ModTime {
					diffs = append(diffs, Diff{
						AbsPath: newChild.AbsPath,
						Op:      2,
					})
				}
			} else if oldChild.IsDir || newChild.IsDir {
				childDiffs := diffNodes(oldChild, newChild)
				diffs = append(diffs, childDiffs...)
			}
		} else {
			diffs = append(diffs, handleNode(oldChild, 0)...)
		}
	}
	for name, newChild := range newChildren {
		_, ok := oldChildren[name]
		if !ok {
			diffs = append(diffs, handleNode(newChild, 1)...)
		}
	}
	return diffs
}
