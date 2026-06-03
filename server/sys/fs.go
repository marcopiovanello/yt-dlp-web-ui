package sys

import (
	"os"
	"path/filepath"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"golang.org/x/sys/unix"
)

// package containing fs related operation (unix only)

// FreeSpace gets the available Bytes writable to download directory
func FreeSpace() (uint64, error) {
	var stat unix.Statfs_t
	unix.Statfs(config.Instance().DownloadPath, &stat)
	return (stat.Bavail * uint64(stat.Bsize)), nil
}

// Build a directory tree started from the specified path using DFS.
// Then return the flattened tree represented as a list.
func DirectoryTree() (*[]string, error) {
	type Node struct {
		path     string
		children []Node
	}

	var (
		rootPath = config.Instance().DownloadPath

		stack     = internal.NewStack[Node]()
		flattened = make([]string, 0)
	)

	stack.Push(Node{path: rootPath})

	flattened = append(flattened, rootPath)

	for stack.IsNotEmpty() {
		current := stack.Pop().Value

		children, err := os.ReadDir(current.path)
		if err != nil {
			return nil, err
		}
		for _, entry := range children {
			var (
				childPath = filepath.Join(current.path, entry.Name())
				childNode = Node{path: childPath}
			)
			if entry.IsDir() {
				current.children = append(current.children, childNode)
				stack.Push(childNode)
				flattened = append(flattened, childNode.path)
			}
		}
	}
	return &flattened, nil
}
