package downloaders

import (
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/common"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
)

type Downloader interface {
	Start() error
	Stop() error
	Status() *internal.ProcessSnapshot

	SetOutput(output internal.DownloadOutput)
	SetProgress(progress internal.DownloadProgress)
	SetMetadata(fetcher func(url string) (*common.DownloadMetadata, error))
	SetPending(p bool)

	IsCompleted() bool

	UpdateSavedFilePath(path string)

	RestoreFromSnapshot(*internal.ProcessSnapshot) error

	GetId() string
	GetUrl() string
}
