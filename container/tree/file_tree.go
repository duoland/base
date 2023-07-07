package tree

import (
	"path/filepath"
	"strings"
)

type FileTreeNode struct {
	Path       string         `json:"-"`
	Title      string         `json:"title"`
	Key        string         `json:"key"`
	Selectable bool           `json:"selectable"` // set to true when file
	IsLeaf     bool           `json:"isLeaf"`     // set to true when file
	Children   *FileTreeNodes `json:"children,omitempty"`
}

type FileTreeNodes []FileTreeNode

func (t *FileTreeNodes) Append(node FileTreeNode) {
	*t = append(*t, node)
}

func CreateFileTreeLayout(filePaths []string) (rootTreeNodes FileTreeNodes) {
	rootTreeNodes = make(FileTreeNodes, 0)
	for _, path := range filePaths {
		var current = &rootTreeNodes
		pathItems := strings.Split(path, "/")
		pathItemsCnt := len(pathItems)
		for index, pathItem := range pathItems {
			var exists bool
			for i := 0; i < len(*current); i++ {
				if (*current)[i].Path == pathItem {
					exists = true
					current = (*current)[i].Children
					break
				}
			}
			if !exists {
				children := new(FileTreeNodes)
				isLeaf := pathItemsCnt-1 == index
				childNode := FileTreeNode{
					Path:       pathItem,
					Title:      pathItem,
					Key:        filepath.Join(pathItems[:index+1]...),
					Children:   children,
					IsLeaf:     isLeaf,
					Selectable: isLeaf,
				}
				current.Append(childNode)
				current = childNode.Children
			}
		}
	}
	return
}
