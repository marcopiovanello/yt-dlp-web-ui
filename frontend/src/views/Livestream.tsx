import {
  Box,
  Button,
  Chip,
  Container,
  Paper,
  Table, TableBody, TableCell, TableContainer, TableHead, TableRow
} from '@mui/material'
import { useAtomValue } from 'jotai'
import { useState } from 'react'
import { interval } from 'rxjs'
import { rpcPollingTimeState } from '../atoms/rpc'
import LivestreamDialog from '../components/livestream/LivestreamDialog'
import LivestreamSpeedDial from '../components/livestream/LivestreamSpeedDial'
import NoLivestreams from '../components/livestream/NoLivestreams'
import LoadingBackdrop from '../components/LoadingBackdrop'
import { useSubscription } from '../hooks/observable'
import { useI18n } from '../hooks/useI18n'
import { useRPC } from '../hooks/useRPC'
import { LiveStreamProgress, LiveStreamStatus } from '../types'

const LiveStreamMonitorView: React.FC = () => {
  const { i18n } = useI18n()
  const { client } = useRPC()

  const [progress, setProgress] = useState<LiveStreamProgress>()
  const [openDialog, setOpenDialog] = useState(false)

  const rpcPollingRate = useAtomValue(rpcPollingTimeState)

  useSubscription(interval(rpcPollingRate), () => {
    client
      .progressLivestream()
      .then(r => setProgress(r.result))
  })

  const formatMicro = (microseconds: number) => {
    const ms = microseconds / 1_000_000
    let s = ms / 1000

    const hr = s / 3600
    s %= 3600

    const mt = s / 60
    s %= 60

    //           huh?
    const ss = (Math.abs(s - 1)).toFixed(0).padStart(2, '0')
    const mts = mt.toFixed(0).padStart(2, '0')
    const hrs = hr.toFixed(0).padStart(2, '0')

    return `${hrs}:${mts}:${ss}`
  }

  const mapStatusToChip = (status: LiveStreamStatus): React.ReactNode => {
    switch (status) {
      case LiveStreamStatus.WAITING:
        return <Chip label='Waiting/Wait start' color='warning' size='small' />
      case LiveStreamStatus.IN_PROGRESS:
        return <Chip label='Downloading' color='primary' size='small' />
      case LiveStreamStatus.COMPLETED:
        return <Chip label='Completed' color='success' size='small' />
      case LiveStreamStatus.ERRORED:
        return <Chip label='Errored' color='error' size='small' />
      default:
        return <Chip label='Unknown state' color='secondary' size='small' />
    }
  }

  const stopAll = () => client.killAllLivestream()
  const stop = (url: string) => client.killLivestream(url)

  return (
    <>
      <LoadingBackdrop isLoading={!progress} />

      <LivestreamSpeedDial onOpen={() => setOpenDialog(s => !s)} onStopAll={stopAll} />
      <LivestreamDialog open={openDialog} onClose={() => setOpenDialog(s => !s)} />

      {!progress || Object.keys(progress).length === 0 ?
        <NoLivestreams /> :
        <Container maxWidth="xl" sx={{ mt: 4, mb: 8 }}>
          <Paper sx={{
            p: 2.5,
            display: 'flex',
            flexDirection: 'column',
            minHeight: '80vh',
          }}>
            <TableContainer component={Box}>
              <Table sx={{ minWidth: '100%' }}>
                <TableHead>
                  <TableRow>
                    <TableCell>{i18n.t('livestreamURLInput')}</TableCell>
                    <TableCell align="right">Status</TableCell>
                    <TableCell align="right">Time to live</TableCell>
                    <TableCell align="right">Starts on</TableCell>
                    <TableCell align="right">Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {progress && Object.keys(progress).map(k => (
                    <TableRow
                      key={k}
                      sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                    >
                      <TableCell>{k}</TableCell>
                      <TableCell align='right'>
                        {mapStatusToChip(progress[k].status)}
                      </TableCell>
                      <TableCell align='right'>
                        {progress[k].status === LiveStreamStatus.WAITING
                          ? formatMicro(Number(progress[k].waitTime))
                          : "-"
                        }
                      </TableCell>
                      <TableCell align='right'>
                        {progress[k].status === LiveStreamStatus.WAITING
                          ? new Date(progress[k].liveDate).toLocaleString()
                          : "-"
                        }
                      </TableCell>
                      <TableCell align='right'>
                        <Button variant='contained' size='small' onClick={() => stop(k)}>
                          Stop
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </Paper>
        </Container>}
    </>
  )
}

export default LiveStreamMonitorView