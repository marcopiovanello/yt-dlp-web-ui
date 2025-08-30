package downloaders

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"regexp"
	"slices"
	"strings"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
)

func argsSanitizer(params []string) []string {
	params = slices.DeleteFunc(params, func(e string) bool {
		match, _ := regexp.MatchString(`(\$\{)|(\&\&)`, e)
		return match
	})

	params = slices.DeleteFunc(params, func(e string) bool {
		return e == ""
	})

	return params
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

func produceLogs(r io.Reader, logs chan<- []byte) {
	go func() {
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			logs <- scanner.Bytes()
		}
	}()
}

func consumeLogs(ctx context.Context, logs <-chan []byte, c LogConsumer, d Downloader) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("detaching logs",
				slog.String("url", d.GetUrl()),
				slog.String("id", c.GetName()),
			)
			return
		case entry := <-logs:
			c.ParseLogEntry(entry, d)
		}
	}
}

func printYtDlpErrors(stdout io.Reader, shortId, url string) {
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		slog.Error("yt-dlp process error",
			slog.String("id", shortId),
			slog.String("url", url),
			slog.String("err", scanner.Text()),
		)
	}
}
