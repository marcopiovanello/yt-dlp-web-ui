import DeleteIcon from '@mui/icons-material/Delete'
import EditIcon from '@mui/icons-material/Edit'
import {
  Box,
  Button,
  Container,
  Paper,
  Table, TableBody, TableCell, TableContainer,
  TableHead, TablePagination, TableRow
} from '@mui/material'
import { matchW } from 'fp-ts/lib/Either'
import { pipe } from 'fp-ts/lib/function'
import { useAtomValue } from 'jotai'
import { useState, useTransition } from 'react'
import { serverURL } from '../atoms/settings'
import LoadingBackdrop from '../components/LoadingBackdrop'
import NoSubscriptions from '../components/subscriptions/NoSubscriptions'
import SubscriptionsDialog from '../components/subscriptions/SubscriptionsDialog'
import SubscriptionsEditDialog from '../components/subscriptions/SubscriptionsEditDialog'
import SubscriptionsSpeedDial from '../components/subscriptions/SubscriptionsSpeedDial'
import { useToast } from '../hooks/toast'
import useFetch from '../hooks/useFetch'
import { useI18n } from '../hooks/useI18n'
import { ffetch } from '../lib/httpClient'
import { Subscription } from '../services/subscriptions'
import { PaginatedResponse } from '../types'

const SubscriptionsView: React.FC = () => {
  const { i18n } = useI18n()
  const { pushMessage } = useToast()

  const baseURL = useAtomValue(serverURL)

  const [selectedSubscription, setSelectedSubscription] = useState<Subscription>()
  const [openDialog, setOpenDialog] = useState(false)

  const [startId, setStartId] = useState(0)
  const [limit, setLimit] = useState(9)
  const [page, setPage] = useState(0)

  const { data: subs, fetcher: refecth } = useFetch<PaginatedResponse<Subscription[]>>(
    `/subscriptions?id=${startId}&limit=${limit}`
  )

  const [isPending, startTransition] = useTransition()

  const deleteSubscription = async (id: string) => {
    const task = ffetch<void>(`${baseURL}/subscriptions/${id}`, {
      method: 'DELETE',
    })
    const either = await task()

    pipe(
      either,
      matchW(
        (l) => pushMessage(l, 'error'),
        () => refecth()
      )
    )
  }

  return (
    <>
      <LoadingBackdrop isLoading={!subs || isPending} />

      <SubscriptionsSpeedDial onOpen={() => setOpenDialog(s => !s)} />

      <SubscriptionsEditDialog
        subscription={selectedSubscription}
        onClose={() => {
          setSelectedSubscription(undefined)
          refecth()
        }}
      />
      <SubscriptionsDialog open={openDialog} onClose={() => {
        setOpenDialog(s => !s)
        refecth()
      }} />

      {!subs || subs.data.length === 0 ?
        <NoSubscriptions /> :
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
                    <TableCell align="left">URL</TableCell>
                    <TableCell align="right">Params</TableCell>
                    <TableCell align="right">{i18n.t('cronExpressionLabel')}</TableCell>
                    <TableCell align="center">Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody sx={{ mb: 'auto' }}>
                  {subs.data.map(x => (
                    <TableRow
                      key={x.id}
                      sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                    >
                      <TableCell>{x.url}</TableCell>
                      <TableCell align='right'>
                        {x.params}
                      </TableCell>
                      <TableCell align='right'>
                        {x.cron_expression}
                      </TableCell>
                      <TableCell align='center'>
                        <Button
                          variant='contained'
                          size='small'
                          sx={{ mr: 0.5 }}
                          onClick={() => setSelectedSubscription(x)}
                        >
                          <EditIcon />
                        </Button>
                        <Button
                          variant='contained'
                          size='small'
                          onClick={() => startTransition(async () => await deleteSubscription(x.id))}
                        >
                          <DeleteIcon />
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
            <TablePagination
              component="div"
              count={-1}
              page={page}
              onPageChange={(_, p) => {
                if (p < page) {
                  setPage(s => (s - 1 <= 0 ? 0 : s - 1))
                  setStartId(subs.first)
                  return
                }
                setPage(s => s + 1)
                setStartId(subs.next)
              }}
              rowsPerPage={limit}
              rowsPerPageOptions={[9, 10, 25, 50, 100]}
              onRowsPerPageChange={(e) => { setLimit(parseInt(e.target.value)) }}
            />
          </Paper>
        </Container>}
    </>
  )
}

export default SubscriptionsView