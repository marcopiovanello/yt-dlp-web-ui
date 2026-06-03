import DeleteIcon from '@mui/icons-material/Delete'
import DownloadIcon from '@mui/icons-material/Download'
import DownloadDoneIcon from '@mui/icons-material/DownloadDone'
import FileDownloadIcon from '@mui/icons-material/FileDownload'
import SmartDisplayIcon from '@mui/icons-material/SmartDisplay'
import StopCircleIcon from '@mui/icons-material/StopCircle'
import {
  Box,
  ButtonGroup,
  IconButton,
  LinearProgress,
  LinearProgressProps,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography
} from "@mui/material"
import { forwardRef } from 'react'
import { TableComponents, TableVirtuoso } from 'react-virtuoso'
import { activeDownloadsState } from '../atoms/downloads'
import { serverURL } from '../atoms/settings'
import { useRPC } from '../hooks/useRPC'
import { ProcessStatus, RPCResult } from '../types'
import { base64URLEncode, formatSize, formatSpeedMiB } from "../utils"
import { useAtomValue } from 'jotai'

const columns = [
  {
    width: 8,
    label: 'Status',
    dataKey: 'status',
    numeric: false,
  },
  {
    width: 500,
    label: 'Title',
    dataKey: 'title',
    numeric: false,
  },
  {
    width: 50,
    label: 'Speed',
    dataKey: 'speed',
    numeric: true,
  },
  {
    width: 150,
    label: 'Progress',
    dataKey: 'progress',
    numeric: true,
  },
  {
    width: 80,
    label: 'Size',
    dataKey: 'size',
    numeric: true,
  },
  {
    width: 100,
    label: 'Added on',
    dataKey: 'addedon',
    numeric: true,
  },
  {
    width: 80,
    label: 'Actions',
    dataKey: 'actions',
    numeric: true,
  },
] as const

function LinearProgressWithLabel(props: LinearProgressProps & { value: number }) {
  return (
    <Box sx={{ display: 'flex', alignItems: 'center' }}>
      <Box sx={{ width: '100%', mr: 1 }}>
        <LinearProgress variant="determinate" {...props} />
      </Box>
      <Box sx={{ minWidth: 35 }}>
        <Typography variant="body2" color="text.secondary">{`${Math.round(
          props.value,
        )}%`}</Typography>
      </Box>
    </Box>
  )
}

const VirtuosoTableComponents: TableComponents<RPCResult> = {
  Scroller: forwardRef<HTMLDivElement>((props, ref) => (
    <TableContainer {...props} ref={ref} />
  )),
  Table: (props) => (
    <Table {...props} sx={{ borderCollapse: 'separate', tableLayout: 'fixed', mt: 2 }} size='small' />
  ),
  TableHead,
  TableRow: ({ item: _item, ...props }) => <TableRow {...props} />,
  TableBody: forwardRef<HTMLTableSectionElement>((props, ref) => (
    <TableBody {...props} ref={ref} />
  )),
}

function fixedHeaderContent() {
  return (
    <TableRow>
      {columns.map((column) => (
        <TableCell
          key={column.dataKey}
          variant="head"
          align={column.numeric || false ? 'right' : 'left'}
          style={{ width: column.width }}
        >
          {column.label}
        </TableCell>
      ))}
    </TableRow>
  )
}

const DownloadsTableView: React.FC = () => {
  const downloads = useAtomValue(activeDownloadsState)
  const serverAddr = useAtomValue(serverURL)
  const { client } = useRPC()

  const viewFile = (path: string) => {
    const encoded = base64URLEncode(path)
    window.open(`${serverAddr}/filebrowser/v/${encoded}?token=${localStorage.getItem('token')}`)
  }

  const downloadFile = (path: string) => {
    const encoded = base64URLEncode(path)
    window.open(`${serverAddr}/filebrowser/d/${encoded}?token=${localStorage.getItem('token')}`)
  }

  const stop = (r: RPCResult) => r.progress.process_status === ProcessStatus.COMPLETED
    ? client.clear(r.id)
    : client.kill(r.id)


  function rowContent(_index: number, download: RPCResult) {
    return (
      <>
        <TableCell>
          {download.progress.percentage === '-1'
            ? <DownloadDoneIcon color="primary" />
            : <DownloadIcon color="primary" />
          }
        </TableCell>
        <TableCell>{download.info.title}</TableCell>
        <TableCell align="right">{formatSpeedMiB(download.progress.speed)}</TableCell>
        <TableCell align="right">
          <LinearProgressWithLabel
            sx={{ height: '16px' }}
            value={
              download.progress.percentage === '-1'
                ? 100
                : Number(download.progress.percentage.replace('%', ''))
            }
            variant={
              download.progress.process_status === 0
                ? 'indeterminate'
                : 'determinate'
            }
            color={download.progress.percentage === '-1' ? 'primary' : 'primary'}
          />
        </TableCell>
        <TableCell align="right">{formatSize(download.info.filesize_approx ?? 0)}</TableCell>
        <TableCell align="right">
          {new Date(download.info.created_at).toLocaleString()}
        </TableCell>
        <TableCell align="right">
          <ButtonGroup>
            <IconButton
              size="small"
              onClick={() => stop(download)}
            >
              {download.progress.percentage === '-1' ? <DeleteIcon /> : <StopCircleIcon />}

            </IconButton>
            {download.progress.percentage === '-1' &&
              <>
                <IconButton
                  size="small"
                  onClick={() => viewFile(download.output.savedFilePath)}
                >
                  <SmartDisplayIcon />
                </IconButton>
                <IconButton
                  size="small"
                  onClick={() => downloadFile(download.output.savedFilePath)}
                >
                  <FileDownloadIcon />
                </IconButton>
              </>
            }
          </ButtonGroup>
        </TableCell>
      </>
    )
  }

  return (
    <Box style={{ height: downloads.length === 0 ? '0vh' : '80vh', width: '100%' }}>
      <TableVirtuoso
        hidden={downloads.length === 0}
        data={downloads}
        components={VirtuosoTableComponents}
        fixedHeaderContent={fixedHeaderContent}
        itemContent={rowContent}
      />
    </Box>
  )
}

export default DownloadsTableView