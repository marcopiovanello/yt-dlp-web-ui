package updater

import (
	"log/slog"
	"os/exec"
	// "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

// Update using the builtin function of yt-dlp
func UpdateExecutable() error {
	// cmd := exec.Command(config.Instance().DownloaderPath, "-U")
	cmd := exec.Command("pip", "install", "-U", "yt-dlp", "--break-system-packages")

	// err := cmd.Start()

	// if err != nil {
	// 	return err
	// }

	// return cmd.Wait()

	out, err := cmd.CombinedOutput()
	slog.Warn(string(out))

	return err

}
