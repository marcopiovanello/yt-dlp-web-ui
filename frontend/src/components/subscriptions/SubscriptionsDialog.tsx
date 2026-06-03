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
  open: boolean
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

const SubscriptionsDialog: React.FC<Props> = ({ open, onClose }) => {
  const [subscriptionURL, setSubscriptionURL] = useState('')
  const [subscriptionCron, setSubscriptionCron] = useState('')

  const customArgs = useAtomValue(customArgsState)

  const { i18n } = useI18n()
  const { pushMessage } = useToast()

  const baseURL = useAtomValue(serverURL)

  const submit = async (sub: Omit<Subscription, 'id'>) => {
    const task = ffetch<void>(`${baseURL}/subscriptions`, {
      method: 'POST',
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
      open={open}
      onClose={onClose}
      TransitionComponent={Transition}
    >
      <AppBar sx={{ position: 'relative' }}>
        <Toolbar>
          <IconButton
            edge="start"
            color="inherit"
            onClick={onClose}
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
                      {i18n.t('subscriptionsInfo')}
                    </Alert>
                    <Alert severity="warning" sx={{ mt: 1 }}>
                      {i18n.t('livestreamExperimentalWarning')}
                    </Alert>
                  </Grid>
                  <Grid item xs={12} mt={1}>
                    <TextField
                      multiline
                      fullWidth
                      label={i18n.t('subscriptionsURLInput')}
                      variant="outlined"
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
                      onChange={(e) => setSubscriptionCron(e.target.value)}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <Button
                      sx={{ mt: 2 }}
                      variant="contained"
                      disabled={subscriptionURL === ''}
                      onClick={() => startTransition(() => submit({
                        url: subscriptionURL,
                        params: customArgs,
                        cron_expression: subscriptionCron
                      }))}
                    >
                      {i18n.t('startButton')}
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

export default SubscriptionsDialog