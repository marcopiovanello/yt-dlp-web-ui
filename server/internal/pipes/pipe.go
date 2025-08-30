package pipes

import "io"

type Pipe interface {
	Name() string
	Connect(r io.Reader) (io.Reader, error)
}
