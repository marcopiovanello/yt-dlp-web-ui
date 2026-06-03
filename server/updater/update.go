package updater

import (
	"os/exec"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

// Update using the builtin function of yt-dlp
func UpdateExecutable() error {
	cmd := exec.Command(config.Instance().DownloaderPath, "-U")

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}
