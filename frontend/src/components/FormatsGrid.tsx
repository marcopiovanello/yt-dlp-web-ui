import { Button, ButtonGroup, Grid, Paper, Typography } from "@mui/material"
import type { DLMetadata } from '../types'

type Props = {
  downloadFormats: DLMetadata
  onAudioSelected: (format: string) => void
  onVideoSelected: (format: string) => void
  onBestQualitySelected: (format: string) => void
  onSubmit: () => void
  onClear: () => void
  pickedBestFormat: string
  pickedAudioFormat: string
  pickedVideoFormat: string
}

export default function FormatsGrid({
  downloadFormats,
  onAudioSelected,
  onVideoSelected,
  onBestQualitySelected,
  onSubmit,
  onClear,
  pickedBestFormat,
  pickedAudioFormat,
  pickedVideoFormat,
}: Props) {
  return (
    <Grid container spacing={2} mt={2}>
      <Grid item xs={12}>
        <Paper
          sx={{
            p: 2,
            display: 'flex',
            flexDirection: 'column',
          }}
        >
          <Grid container>
            <Grid item xs={12}>
              <Typography variant="h6" component="div" pb={1}>
                {downloadFormats.title}
              </Typography>
              {/* <Skeleton variant="rectangular" height={180} /> */}
            </Grid>
            <Grid item xs={12} pb={1}>
              <img src={downloadFormats.thumbnail} height={260} width="100%" style={{ objectFit: 'cover' }} />
            </Grid>
            {/* video only */}
            <Grid item xs={12}>
              <Typography variant="body1" component="div">
                Best quality
              </Typography>
            </Grid>
            <Grid item pr={2} py={1}>
              <Button
                variant="contained"
                disabled={pickedBestFormat !== ''}
                onClick={() => onBestQualitySelected(downloadFormats.best.format_id)}
              >
                {downloadFormats.best.format_note || downloadFormats.best.format_id} - {downloadFormats.best.vcodec}+{downloadFormats.best.acodec}
                &nbsp;({downloadFormats.best.resolution}{(downloadFormats.best.filesize_approx > 0) ? ", ~" + Math.round(downloadFormats.best.filesize_approx / 1024 / 1024) + " MiB" : ""})
              </Button>
            </Grid>
            {/* video only */}
            {downloadFormats.formats.filter(format => format.acodec === 'none' && format.vcodec !== 'none').length &&
              <Grid item xs={12}>
                <Typography variant="body1" component="div">
                  Video data {downloadFormats.formats[1].acodec}
                </Typography>
              </Grid>
            }
            {downloadFormats.formats
              .filter(format => format.acodec === 'none' && format.vcodec !== 'none')
              .map((format, idx) => (
                <Grid item pr={2} py={1} key={idx}>
                  <Button
                    variant="contained"
                    onClick={() => onVideoSelected(format.format_id)}
                    disabled={pickedVideoFormat === format.format_id}
                  >
                    {format.format_note} - {format.vcodec === 'none' ? format.acodec : format.vcodec}
                    &nbsp;({format.resolution}{(format.filesize_approx > 0) ? ", ~" + Math.round(format.filesize_approx / 1024 / 1024) + " MiB" : ""})
                  </Button>
                </Grid>
              ))
            }
            {downloadFormats.formats.filter(format => format.acodec === 'none' && format.vcodec !== 'none').length &&
              <Grid item xs={12}>
                <Typography variant="body1" component="div">
                  Audio data
                </Typography>
              </Grid>
            }
            {downloadFormats.formats
              .filter(format => format.acodec !== 'none' && format.vcodec === 'none')
              .map((format, idx) => (
                <Grid item pr={2} py={1} key={idx}>
                  <Button
                    variant="contained"
                    onClick={() => onAudioSelected(format.format_id)}
                    disabled={pickedAudioFormat === format.format_id}
                  >
                    {format.format_note} - {format.vcodec === 'none' ? format.acodec : format.vcodec}
                    {(format.filesize_approx > 0) ? " (~" + Math.round(format.filesize_approx / 1024 / 1024) + " MiB)" : ""}
                    {format.language}
                  </Button>
                </Grid>
              ))
            }
            <Grid item xs={12} pt={2}>
              <ButtonGroup disableElevation variant="contained">
                <Button
                  onClick={() => onSubmit()}
                  disabled={!pickedBestFormat && !(pickedAudioFormat || pickedVideoFormat)}
                > Download
                </Button>
                <Button
                  onClick={() => onClear()}
                > Clear
                </Button>
              </ButtonGroup>
            </Grid>
          </Grid>
        </Paper>
      </Grid>
    </Grid>
  )
}