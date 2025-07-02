package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileNode struct {
	Name     string
	Path     string
	IsDir    bool
	IsOpen   bool
	Selected bool
	Children []*FileNode
	Parent   *FileNode
}

func BuildFileTree(root string) *FileNode {
	rootNode := &FileNode{
		Name:   filepath.Base(root),
		Path:   root,
		IsDir:  true,
		IsOpen: true,
	}
	buildFileTreeRecursive(rootNode, root, 0)
	return rootNode
}

func buildFileTreeRecursive(parent *FileNode, path string, depth int) {
	if depth > 3 {
		return
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}
	// Sort entries: directories first, then files
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") ||
			entry.Name() == "node_modules" ||
			entry.Name() == "vendor" {
			continue
		}
		childPath := filepath.Join(path, entry.Name())
		child := &FileNode{
			Name:   entry.Name(),
			Path:   childPath,
			IsDir:  entry.IsDir(),
			IsOpen: false,
			Parent: parent,
		}
		parent.Children = append(parent.Children, child)
		if entry.IsDir() {
			buildFileTreeRecursive(child, childPath, depth+1)
		}
	}
}

func FlattenFileTree(root *FileNode) []*FileNode {
	var result []*FileNode
	flattenFileTreeRecursive(root, &result, 0)
	return result
}

func flattenFileTreeRecursive(node *FileNode, result *[]*FileNode, depth int) {
	if depth > 0 {
		*result = append(*result, node)
	}
	if node.IsDir && node.IsOpen {
		for _, child := range node.Children {
			flattenFileTreeRecursive(child, result, depth+1)
		}
	}
}

func RenderFileNode(node *FileNode) string {
	depth := GetNodeDepth(node)
	indent := strings.Repeat("  ", depth)
	if node.IsDir {
		icon := "▶"
		if node.IsOpen {
			icon = "▼"
		}
		fileCount := ""
		if !node.IsOpen && len(node.Children) > 0 {
			fileCount = fmt.Sprintf(" (%d items)", len(node.Children))
		}
		return fmt.Sprintf("%s%s %s/%s", indent, icon, node.Name, fileCount)
	} else {
		checkbox := "◯"
		if node.Selected {
			checkbox = "◉"
		}
		return fmt.Sprintf("%s  %s %s", indent, checkbox, node.Name)
	}
}

func GetNodeDepth(node *FileNode) int {
	depth := 0
	current := node
	for current.Parent != nil {
		depth++
		current = current.Parent
	}
	return depth
}
