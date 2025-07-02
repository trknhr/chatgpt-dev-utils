package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildFileTreeAndFlatten(t *testing.T) {
	t.Run("basic directory structure", func(t *testing.T) {
		dir := t.TempDir()

		err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0644)
		require.NoError(t, err)

		err = os.Mkdir(filepath.Join(dir, "subdir"), 0755)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(dir, "subdir", "b.txt"), []byte("world"), 0644)
		require.NoError(t, err)

		tree := BuildFileTree(dir)
		require.NotNil(t, tree, "tree should not be nil")
		assert.True(t, tree.IsDir, "root node should be a directory")

		flat := FlattenFileTree(tree)
		assert.Len(t, flat, 2, "expected 2 top-level entries (a.txt, subdir)")
	})
}

func TestRenderFileNode(t *testing.T) {
	t.Run("render selected file node", func(t *testing.T) {
		n := &FileNode{Name: "foo.txt", Path: "foo.txt", IsDir: false, Selected: true}
		out := RenderFileNode(n)
		assert.NotEmpty(t, out, "render output should not be empty")
	})
}

func TestGetNodeDepth(t *testing.T) {
	t.Run("grandchild depth", func(t *testing.T) {
		root := &FileNode{Name: "root", IsDir: true}
		child := &FileNode{Name: "child", Parent: root}
		grandchild := &FileNode{Name: "grandchild", Parent: child}

		depth := GetNodeDepth(grandchild)
		assert.Equal(t, 2, depth, "grandchild depth should be 2")
	})
}
