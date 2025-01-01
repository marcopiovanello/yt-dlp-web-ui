package livestream

import (
	"testing"
	"time"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
)

func setupTest() {
	config.Instance().DownloaderPath = "yt-dlp"
}

func TestLivestream(t *testing.T) {
	setupTest()

	done := make(chan *LiveStream)

	ls := New("https://www.youtube.com/watch?v=LSm1daKezcE", done, &internal.MessageQueue{}, &internal.MemoryDB{})
	go ls.Start()

	time.AfterFunc(time.Second*20, func() {
		ls.Kill()
	})

	for {
		select {
		case wt := <-ls.WaitTime():
			t.Log(wt)
		case <-done:
			t.Log("done")
			return
		}
	}
}
