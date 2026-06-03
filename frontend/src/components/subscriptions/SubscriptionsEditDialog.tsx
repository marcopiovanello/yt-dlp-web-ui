import CloseIcon from '@mui/icons-material/Close'
import {
  Alert,
  AppBar,
  Box,
  Button,
  Container,
  Dialog,
  Grid,
  IconButton,
  Paper,
  Slide,
  TextField,
  Toolbar,
  Typography
} from '@mui/material'
import { TransitionProps } from '@mui/material/transitions'
import { matchW } from 'fp-ts/lib/Either'
import { pipe } from 'fp-ts/lib/function'
import { useAtomValue } from 'jotai'
import { forwardRef, startTransition, useState } from 'react'
import { customArgsState } from '../../atoms/downloadTemplate'
import { serverURL } from '../../atoms/settings'
import { useToast } from '../../hooks/toast'
import { useI18n } from '../../hooks/useI18n'
import { ffetch } from '../../lib/httpClient'
import { Subscription } from '../../services/subscriptions'
import ExtraDownloadOptions from '../ExtraDownloadOptions'

type Props = {
  subscription: Subscription | undefined
  onClose: () => void
}

const Transition = forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement
  },
  ref: React.Ref<unknown>,
) {
  return <Slide direction="up" ref={ref} {...props} />
})

const SubscriptionsEditDialog: React.FC<Props> = ({ subscription, onClose }) => {
  const [subscriptionURL, setSubscriptionURL] = useState('')
  const [subscriptionCron, setSubscriptionCron] = useState('')

  const customArgs = useAtomValue(customArgsState)

  const { i18n } = useI18n()
  const { pushMessage } = useToast()

  const baseURL = useAtomValue(serverURL)

  const editSubscription = async (sub: Subscription) => {
    const task = ffetch<void>(`${baseURL}/subscriptions`, {
      method: 'PATCH',
      body: JSON.stringify(sub)
    })
    const either = await task()

    pipe(
      either,
      matchW(
        (l) => pushMessage(l, 'error'),
        (_) => onClose()
      )
    )
  }

  return (
    <Dialog
      fullScreen
      open={!!subscription}
      TransitionComponent={Transition}
    >
      <AppBar sx={{ position: 'relative' }}>
        <Toolbar>
          <IconButton
            edge="start"
            color="inherit"
            onClick={() => onClose()}
            aria-label="close"
          >
            <CloseIcon />
          </IconButton>
          <Typography sx={{ ml: 2, flex: 1 }} variant="h6" component="div">
            {i18n.t('subscriptionsButtonLabel')}
          </Typography>
        </Toolbar>
      </AppBar>
      <Box sx={{
        backgroundColor: (theme) => theme.palette.background.default,
        minHeight: (theme) => `calc(99vh - ${theme.mixins.toolbar.minHeight}px)`
      }}>
        <Container sx={{ my: 4 }}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <Paper
                elevation={4}
                sx={{
                  p: 2,
                  display: 'flex',
                  flexDirection: 'column',
                }}
              >
                <Grid container gap={1.5}>
                  <Grid item xs={12}>
                    <Alert severity="info">
                      Editing {subscription?.url}
                    </Alert>
                  </Grid>
                  <Grid item xs={12} mt={1}>
                    <TextField
                      multiline
                      fullWidth
                      label={i18n.t('subscriptionsURLInput')}
                      variant="outlined"
                      defaultValue={subscription?.url}
                      placeholder="https://www.youtube.com/@SomeChannelThatExists/videos"
                      onChange={(e) => setSubscriptionURL(e.target.value)}
                    />
                  </Grid>
                  <Grid item xs={8} mt={-2}>
                    <ExtraDownloadOptions />
                  </Grid>
                  <Grid item xs={3.871}>
                    <TextField
                      multiline
                      fullWidth
                      label={i18n.t('cronExpressionLabel')}
                      variant="outlined"
                      placeholder="*/5 * * * *"
                      defaultValue={subscription?.cron_expression}
                      onChange={(e) => setSubscriptionCron(e.target.value)}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <Button
                      sx={{ mt: 2 }}
                      variant="contained"
                      onClick={() => startTransition(async () => await editSubscription({
                        id: subscription?.id ?? '',
                        url: subscriptionURL || subscription?.url!,
                        params: customArgs || subscription?.params!,
                        cron_expression: subscriptionCron || subscription?.cron_expression!
                      }))}
                    >
                      {i18n.t('editButtonLabel')}
                    </Button>
                  </Grid>
                </Grid>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>
    </Dialog>
  )
}

export default SubscriptionsEditDialog