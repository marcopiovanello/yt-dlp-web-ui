package logging

import (
	"compress/gzip"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

/*
implements io.Writer interface

File base logger with log-rotate capabilities.
The rotate process must be initiated from an external goroutine.

After rotation the previous logs file are compressed with gzip algorithm.

The rotated log follows this naming: [filename].UTC time.gz
*/
type LogRotateWriter struct {
	mu       sync.Mutex
	fd       *os.File
	filename string
}

func NewRotableLogger(filename string) (*LogRotateWriter, error) {
	fd, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	w := &LogRotateWriter{filename: filename, fd: fd}
	return w, nil
}

func (w *LogRotateWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.fd.Write(b)
}

func (w *LogRotateWriter) Rotate() error {
	slog.Info("started log rotation")

	w.mu.Lock()

	gzFile, err := os.Create(strings.TrimSuffix(w.filename, ".log") + "-" + time.Now().Format(time.RFC3339) + ".log.gz")
	if err != nil {
		return err
	}

	zw := gzip.NewWriter(gzFile)

	defer func() {
		zw.Close()
		zw.Flush()
		gzFile.Close()
	}()

	if _, err := os.Stat(w.filename); err != nil {
		return err
	}

	fd, _ := os.Open(w.filename)
	io.Copy(zw, fd)
	fd.Close()

	w.fd.Close()

	if err := os.Remove(w.filename); err != nil {
		return err
	}

	w.fd, _ = os.Create(w.filename)

	w.mu.Unlock()
	slog.Info("ended log rotation")

	return err
}
