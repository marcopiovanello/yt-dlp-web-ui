package playlist

import (
	"slices"
	"strconv"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/common"
)

/*
	Applicable modifiers

				full									|		   short		  |					description
	---------------------------------------------------------------------------------
	--playlist-start NUMBER     |    -I NUMBER:	  |	  discard first N entries
	--playlist-end NUMBER       |    -I :NUMBER   |   discard last N entries
	--playlist-reverse          |    -I ::-1			|   self explanatory
	--max-downloads NUMBER      |                 |   stops after N completed downloads
*/

func ApplyModifiers(entries *[]common.DownloadInfo, args []string) error {
	for i, modifier := range args {
		switch modifier {
		case "--playlist-start":
			return playlistStart(i, modifier, args, entries)

		case "--playlist-end":
			return playlistEnd(i, modifier, args, entries)

		case "--max-downloads":
			return maxDownloads(i, modifier, args, entries)

		case "--playlist-reverse":
			slices.Reverse(*entries)
			return nil
		}
	}
	return nil
}

func playlistStart(i int, modifier string, args []string, entries *[]common.DownloadInfo) error {
	if !guard(i, len(modifier)) {
		return nil
	}

	n, err := strconv.Atoi(args[i+1])
	if err != nil {
		return err
	}

	*entries = (*entries)[n:]

	return nil
}

func playlistEnd(i int, modifier string, args []string, entries *[]common.DownloadInfo) error {
	if !guard(i, len(modifier)) {
		return nil
	}

	n, err := strconv.Atoi(args[i+1])
	if err != nil {
		return err
	}

	*entries = (*entries)[:n]

	return nil
}

func maxDownloads(i int, modifier string, args []string, entries *[]common.DownloadInfo) error {
	if !guard(i, len(modifier)) {
		return nil
	}

	n, err := strconv.Atoi(args[i+1])
	if err != nil {
		return err
	}

	*entries = (*entries)[0:n]

	return nil
}

func guard(i, len int) bool { return i+1 < len-1 }
