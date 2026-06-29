package pipes

import (
	"io"
	"log/slog"
	"os"
)

type FileWriter struct {
	Path    string
	IsFinal bool
}

func (f *FileWriter) Name() string { return "file-writer" }

func (f *FileWriter) Connect(r io.Reader) (io.Reader, error) {
	file, err := os.Create(f.Path)
	if err != nil {
		return nil, err
	}

	if f.IsFinal {
		go func() {
			defer file.Close()
			if _, err := io.Copy(file, r); err != nil {
				slog.Error("FileWriter (final) error", slog.Any("err", err))
			}
		}()
		return r, nil
	}

	pr, pw := io.Pipe()

	go func() {
		defer file.Close()
		defer pw.Close()

		writer := io.MultiWriter(file, pw)
		if _, err := io.Copy(writer, r); err != nil {
			slog.Error("FileWriter (pipeline) error", slog.Any("err", err))
		}
	}()

	return pr, nil
}
