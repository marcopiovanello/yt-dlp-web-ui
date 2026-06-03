import DeleteIcon from '@mui/icons-material/Delete'
import DeleteForeverIcon from '@mui/icons-material/DeleteForever'
import OpenInBrowserIcon from '@mui/icons-material/OpenInBrowser'
import SaveAltIcon from '@mui/icons-material/SaveAlt'
import {
  Card,
  CardActionArea,
  CardActions,
  CardContent,
  CardMedia,
  IconButton,
  Skeleton,
  Tooltip,
  Typography
} from '@mui/material'
import { useAtomValue } from 'jotai'
import { serverURL } from '../atoms/settings'
import { ArchiveEntry } from '../types'
import { base64URLEncode, ellipsis } from '../utils'

type Props = {
  entry: ArchiveEntry
  onDelete: (id: string) => void
  onHardDelete: (id: string) => void
}

const ArchiveCard: React.FC<Props> = ({ entry, onDelete, onHardDelete }) => {
  const serverAddr = useAtomValue(serverURL)

  const viewFile = (path: string) => {
    const encoded = base64URLEncode(path)
    window.open(`${serverAddr}/filebrowser/v/${encoded}?token=${localStorage.getItem('token')}`)
  }

  const downloadFile = (path: string) => {
    const encoded = base64URLEncode(path)
    window.open(`${serverAddr}/filebrowser/d/${encoded}?token=${localStorage.getItem('token')}`)
  }

  return (
    <Card>
      <CardActionArea onClick={() => navigator.clipboard.writeText(entry.source)}>
        {entry.thumbnail !== '' ?
          <CardMedia
            component="img"
            height={180}
            image={entry.thumbnail}
          /> :
          <Skeleton variant="rectangular" height={180} />
        }
        <CardContent>
          {entry.title !== '' ?
            <Typography gutterBottom variant="h6" component="div">
              {ellipsis(entry.title, 60)}
            </Typography> :
            <Skeleton />
          }
          {/* <code>
            {JSON.stringify(JSON.parse(entry.metadata), null, 2)}
          </code> */}
          <p>{new Date(entry.created_at).toLocaleString()}</p>
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Tooltip title="Open in browser">
          <IconButton
            onClick={() => viewFile(entry.path)}
          >
            <OpenInBrowserIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Download this file">
          <IconButton
            onClick={() => downloadFile(entry.path)}
          >
            <SaveAltIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete from archive">
          <IconButton
            onClick={() => onDelete(entry.id)}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete from disk">
          <IconButton
            onClick={() => onHardDelete(entry.id)}
          >
            <DeleteForeverIcon />
          </IconButton>
        </Tooltip>
      </CardActions>
    </Card>
  )
}

export default ArchiveCard