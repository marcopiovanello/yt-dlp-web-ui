package updater

import (
	"os/exec"
)

// Update using the builtin function of yt-dlp
func UpdateExecutable() error {
	cmd := exec.Command("pip", "install", "-U", "yt-dlp", "--break-system-packages")

	err := cmd.Start()

	if err != nil {
		return err
	}

	return cmd.Wait()

}
