import {
  Chip,
  Container,
  Paper
} from '@mui/material'
import { matchW } from 'fp-ts/lib/Either'
import { pipe } from 'fp-ts/lib/function'
import { useAtomValue } from 'jotai'
import { useState, useTransition } from 'react'
import { serverURL } from '../atoms/settings'
import LoadingBackdrop from '../components/LoadingBackdrop'
import NoSubscriptions from '../components/subscriptions/NoSubscriptions'
import SubscriptionsSpeedDial from '../components/subscriptions/SubscriptionsSpeedDial'
import TwitchDialog from '../components/twitch/TwitchDialog'
import { useToast } from '../hooks/toast'
import useFetch from '../hooks/useFetch'
import { ffetch } from '../lib/httpClient'

const TwitchView: React.FC = () => {
  const { pushMessage } = useToast()

  const baseURL = useAtomValue(serverURL)

  const [openDialog, setOpenDialog] = useState(false)

  const { data: users, fetcher: refetch } = useFetch<Array<string>>('/twitch/users')

  const [isPending, startTransition] = useTransition()

  const deleteUser = async (user: string) => {
    const task = ffetch<void>(`${baseURL}/twitch/user/${user}`, {
      method: 'DELETE',
    })
    const either = await task()

    pipe(
      either,
      matchW(
        (l) => pushMessage(l, 'error'),
        () => refetch()
      )
    )
  }

  return (
    <>
      <LoadingBackdrop isLoading={!users || isPending} />

      <SubscriptionsSpeedDial onOpen={() => setOpenDialog(s => !s)} />

      <TwitchDialog open={openDialog} onClose={() => {
        setOpenDialog(s => !s)
        refetch()
      }} />

      {
        !users || users.length === 0 ?
          <NoSubscriptions /> :
          <Container maxWidth="xl" sx={{ mt: 4, mb: 8 }}>
            <Paper sx={{
              p: 2.5,
              minHeight: '80vh',
            }}>
              {users.map(user => (
                <Chip
                  label={user}
                  onDelete={() => startTransition(async () => await deleteUser(user))}
                />
              ))}
            </Paper>
          </Container>
      }
    </>
  )
}

export default TwitchView