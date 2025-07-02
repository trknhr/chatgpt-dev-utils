package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFileTreeAndFlatten(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0644)
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "subdir", "b.txt"), []byte("world"), 0644)

	tree := BuildFileTree(dir)
	assert.NotNil(t, tree, "expected root to be non-nil")
	assert.True(t, tree.IsDir, "expected root to be a directory node")
	flat := FlattenFileTree(tree)
	assert.Equal(t, 2, len(flat), "expected 2 files (a.txt, subdir)")
}

func TestRenderFileNode(t *testing.T) {
	n := &FileNode{Name: "foo.txt", Path: "foo.txt", IsDir: false, Selected: true}
	out := RenderFileNode(n)
	assert.NotEmpty(t, out, "expected non-empty render output")
}

func TestGetNodeDepth(t *testing.T) {
	root := &FileNode{Name: "root", IsDir: true}
	child := &FileNode{Name: "child", Parent: root}
	grandchild := &FileNode{Name: "grandchild", Parent: child}
	assert.Equal(t, 2, GetNodeDepth(grandchild), "expected depth 2")
}
