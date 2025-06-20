import {
  Button,
  Card,
  CardActionArea,
  CardActions,
  CardContent,
  CardMedia,
  Chip,
  IconButton,
  LinearProgress,
  Skeleton,
  Stack,
  Tooltip,
  Typography
} from '@mui/material'
import { useAtomValue } from 'jotai'
import { useCallback } from 'react'
import { serverURL } from '../atoms/settings'
import { RPCResult } from '../types'
import { base64URLEncode, ellipsis, formatSize, formatSpeedMiB, mapProcessStatus } from '../utils'
import ResolutionBadge from './ResolutionBadge'
import ClearIcon from '@mui/icons-material/Clear'
import StopCircleIcon from '@mui/icons-material/StopCircle'
import OpenInBrowserIcon from '@mui/icons-material/OpenInBrowser'
import SaveAltIcon from '@mui/icons-material/SaveAlt'

import ShareIcon from '@mui/icons-material/Share'
import ContentCopyIcon from '@mui/icons-material/ContentCopy'
import { useState } from 'react'
import { Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions } from '@mui/material'


type Props = {
  download: RPCResult
  onStop: () => void
  onCopy: () => void
}

const DownloadCard: React.FC<Props> = ({ download, onStop, onCopy }) => {
  const serverAddr = useAtomValue(serverURL)

  const isCompleted = useCallback(
    () => download.progress.percentage === '-1',
    [download.progress.percentage]
  )

  const percentageToNumber = useCallback(
    () => isCompleted()
      ? 100
      : Number(download.progress.percentage.replace('%', '')),
    [download.progress.percentage, isCompleted]
  )

  const viewFile = (path: string) => {
    const encoded = base64URLEncode(path)
    window.open(`${serverAddr}/filebrowser/v/${encoded}?token=${localStorage.getItem('token')}`)
  }

  const downloadFile = (path: string) => {
    const encoded = base64URLEncode(path)
    window.open(`${serverAddr}/filebrowser/d/${encoded}?token=${localStorage.getItem('token')}`)
  }

  const handleShare = () => {
    const encoded = base64URLEncode(download.output.savedFilePath)
    //const link = `${serverAddr}/public/${encoded}`
    const link = `${window.location.origin}/#/public/${encoded}`
    setShareLink(link)
    setShareOpen(true)
  }

  const [shareOpen, setShareOpen] = useState(false)
  const [shareLink, setShareLink] = useState('')


  return (
    <>
      <Card>
        <CardActionArea onClick={() => {
          navigator.clipboard.writeText(download.info.url)
          onCopy()
        }}>
          {download.info.thumbnail !== '' ?
            <CardMedia
              component="img"
              height={180}
              image={download.info.thumbnail}
            /> :
            <Skeleton variant="rectangular" height={180} />
          }
          {download.progress.percentage ?
            <LinearProgress
              variant="determinate"
              value={percentageToNumber()}
              color={isCompleted() ? "success" : "primary"}
            /> :
            null
          }
          <CardContent>
            {download.info.title !== '' ?
              <Typography gutterBottom variant="h6" component="div">
                {ellipsis(download.info.title, 100)}
              </Typography> :
              <Skeleton />
            }
            <Stack direction="row" spacing={0.5} py={1}>
              <Chip
                label={
                  isCompleted()
                    ? 'Completed'
                    : mapProcessStatus(download.progress.process_status)
                }
                color="primary"
                size="small"
              />
              <Typography>
                {!isCompleted() ? download.progress.percentage : ''}
              </Typography>
              <Typography>
                &nbsp;
                {!isCompleted() ? formatSpeedMiB(download.progress.speed) : ''}
              </Typography>
              <Typography>
                {formatSize(download.info.filesize_approx ?? 0)}
              </Typography>
              <ResolutionBadge resolution={download.info.resolution} />
            </Stack>
          </CardContent>
        </CardActionArea>
        <CardActions>
          {isCompleted() ?
            <Tooltip title="Clear from the view">
              <IconButton
                onClick={onStop}
              >
                <ClearIcon />
              </IconButton>
            </Tooltip>
            :
            <Tooltip title="Stop this download">
              <IconButton
                onClick={onStop}
              >
                <StopCircleIcon />
              </IconButton>
            </Tooltip>
          }
          {isCompleted() &&
            <>
              <Tooltip title="Download this file">
                <IconButton
                  onClick={() => downloadFile(download.output.savedFilePath)}
                >
                  <SaveAltIcon />
                </IconButton>
              </Tooltip>
              <Tooltip title="Open in a new tab">
                <IconButton
                  onClick={() => viewFile(download.output.savedFilePath)}
                >
                  <OpenInBrowserIcon />
                </IconButton>
              </Tooltip>
              <Tooltip title="Share this file">  {/* ← Das hier neu */}
                <IconButton onClick={handleShare}>
                  <ShareIcon />
                </IconButton>
              </Tooltip>
            </>
          }
        </CardActions>
      </Card>
      <Dialog
        open={shareOpen}
        onClose={() => setShareOpen(false)}
      >
        <DialogTitle>Share this file</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Copy the link below and share it.
          </DialogContentText>
          <Stack direction="row" spacing={1} alignItems="center" mt={2}>
            <Typography
              sx={{
                fontFamily: 'monospace',
                fontSize: '0.8rem',
                overflowWrap: 'anywhere',
                flex: 1,
              }}
            >
              {shareLink}
            </Typography>
            <Button
              variant="outlined"
              onClick={() => {
                navigator.clipboard.writeText(shareLink)
              }}
            >
              <ContentCopyIcon fontSize="small" sx={{ mr: 1 }} />
              Copy
            </Button>
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShareOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </>
  )
}

export default DownloadCard