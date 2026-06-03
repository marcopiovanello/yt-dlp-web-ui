package archive

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

// Perform a search on the archive.txt file an determines if a download
// has already be done.
func DownloadExists(ctx context.Context, url string) (bool, error) {
	cmd := exec.CommandContext(
		ctx,
		config.Instance().DownloaderPath,
		"--print",
		"%(extractor)s %(id)s",
		url,
	)
	stdout, err := cmd.Output()
	if err != nil {
		return false, err
	}

	extractorAndURL := bytes.Trim(stdout, "\n")

	fd, err := os.Open(filepath.Join(config.Instance().Dir(), "archive.txt"))
	if err != nil {
		return false, err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)

	// search linearly for lower memory usage...
	// the a pre-sorted with hashed values version of the archive.txt file can be loaded in memory
	// and perform a binary search on it.
	for scanner.Scan() {
		if bytes.Equal(scanner.Bytes(), extractorAndURL) {
			return true, nil
		}
	}

	// data, err := io.ReadAll(fd)
	// if err != nil {
	// 	return false, err
	// }

	// slices.BinarySearchFunc(data, extractorAndURL, func(a []byte, b []byte) int {
	// 	return hash(a).Compare(hash(b))
	// })

	return false, nil
}
