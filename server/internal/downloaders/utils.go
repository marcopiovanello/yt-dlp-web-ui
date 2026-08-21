package downloaders

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"slices"
	"strings"

	"github.com/marcopiovanello/yt-dlp-web-ui/v4/server/internal"
)

var allowedFlags = map[string]bool{
	"-f":                               true, // format selection
	"--format":                         true,
	"-S":                               true,
	"--format-sort":                    true,
	"--format-sort-reset":              true,
	"--format-sort-force":              true,
	"--no-format-sort-force":           true,
	"--video-multistreams":             true,
	"--no-video-multistreams":          true,
	"--audio-multistreams":             true,
	"--no-audio-multistreams":          true,
	"--prefer-free-formats":            true,
	"--no-prefer-free-formats":         true,
	"--check-formats":                  true,
	"--check-all-formats":              true,
	"--no-check-formats":               true,
	"-F":                               true,
	"--list-formats":                   true,
	"--merge-output-format":            true,
	"-I":                               true, // video / playlist selection
	"--playlist-items":                 true,
	"--min-filesize":                   true,
	"--max-filesize":                   true,
	"--date":                           true,
	"--datebefore":                     true,
	"--dateafter":                      true,
	"--match-filters":                  true,
	"--no-match-filters":               true,
	"--break-match-filters":            true,
	"--no-break-match-filters":         true,
	"--yes-playlist":                   true,
	"--age-limit":                      true,
	"--max-downloads":                  true,
	"--playlist-random":                true,
	"-N":                               true, // download behaviour
	"--concurrent-fragments":           true,
	"-r":                               true,
	"--limit-rate":                     true,
	"--throttled-rate":                 true,
	"-R":                               true,
	"--retries":                        true,
	"--fragment-retries":               true,
	"--retry-sleep":                    true,
	"--skip-unavailable-fragments":     true,
	"--abort-on-unavailable-fragments": true,
	"--buffer-size":                    true,
	"--resize-buffer":                  true,
	"--no-resize-buffer":               true,
	"--http-chunk-size":                true,
	"--download-sections":              true,
	"--hls-use-mpegts":                 true,
	"--no-hls-use-mpegts":              true,
	"--write-subs":                     true, // subs
	"--no-write-subs":                  true,
	"--write-auto-subs":                true,
	"--no-write-auto-subs":             true,
	"--list-subs":                      true,
	"--sub-format":                     true,
	"--sub-langs":                      true,
	"--convert-subs":                   true,
	"--write-thumbnail":                true, // thumbs
	"--write-all-thumbnails":           true,
	"--convert-thumbnails":             true,
	"--embed-subs":                     true,
	"--embed-thumbnail":                true,
	"--embed-metadata":                 true,
	"--embed-chapters":                 true,
	"--embed-info-json":                true,
	"--write-description":              true, // meta
	"--audio-format":                   true, // safe post processing
	"--audio-quality":                  true,
	"--remux-video":                    true,
	"--recode-video":                   true,
	"--split-chapters":                 true,
	"--no-split-chapters":              true,
	"--remove-chapters":                true,
	"--no-remove-chapters":             true,
	"--force-keyframes-at-cuts":        true,
	"--no-force-keyframes-at-cuts":     true,
	"--fixup":                          true,
	"--sponsorblock-mark":              true, // sponsorblock
	"--sponsorblock-remove":            true,
	"--sponsorblock-chapter-title":     true,
	"--no-sponsorblock":                true,
	"-v":                               true, // verbosity
	"--verbose":                        true,
}

func argsSanitizer(params []string) ([]string, error) {
	params = slices.DeleteFunc(params, func(e string) bool {
		match, _ := regexp.MatchString(`(\$\{)|(\&\&)`, e)
		return match
	})

	params = slices.DeleteFunc(params, func(e string) bool {
		return e == ""
	})

	var out []string

	for i := 0; i < len(params); i++ {
		p := params[i]
		if !strings.HasPrefix(p, "-") {
			out = append(out, p)
			continue
		}
		if !allowedFlags[p] {
			return nil, fmt.Errorf("param %s not allowed", p)
		}
		out = append(out, p)
	}

	return out, nil
}

func buildFilename(o *internal.DownloadOutput) {
	if o.Filename != "" && strings.Contains(o.Filename, ".%(ext)s") {
		o.Filename += ".%(ext)s"
	}

	o.Filename = strings.Replace(
		o.Filename,
		".%(ext)s.%(ext)s",
		".%(ext)s",
		1,
	)
}

func produceLogs(r io.Reader, logs chan<- []byte) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		logs <- scanner.Bytes()
	}

	return scanner.Err()
}

func consumeLogs(ctx context.Context, logs <-chan []byte, c LogConsumer, d Downloader) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("detaching from yt-dlp log buffer",
				slog.String("url", d.GetUrl()),
				slog.String("id", c.GetName()),
			)
			return
		case entry := <-logs:
			c.ParseLogEntry(entry, d)
		}
	}
}

func printYtDlpErrors(stdout io.Reader, shortId, url string) error {
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		slog.Error("yt-dlp process error",
			slog.String("id", shortId),
			slog.String("url", url),
			slog.String("err", scanner.Text()),
		)
	}

	return scanner.Err()
}
