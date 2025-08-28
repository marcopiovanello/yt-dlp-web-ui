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
import { serverURL } from '../../atoms/settings'
import { useToast } from '../../hooks/toast'
import { useI18n } from '../../hooks/useI18n'
import { ffetch } from '../../lib/httpClient'

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

const TwitchDialog: React.FC<Props> = ({ open, onClose }) => {
  const [channelURL, setChannelURL] = useState('')

  const { i18n } = useI18n()
  const { pushMessage } = useToast()

  const baseURL = useAtomValue(serverURL)

  const submit = async (channelURL: string) => {
    const task = ffetch<void>(`${baseURL}/twitch/user`, {
      method: 'POST',
      body: JSON.stringify({
        user: channelURL.split('/').at(-1)
      })
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
                      placeholder="https://www.twitch.tv/a_twitch_user_that_exists"
                      onChange={(e) => setChannelURL(e.target.value)}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <Button
                      sx={{ mt: 2 }}
                      variant="contained"
                      disabled={channelURL === ''}
                      onClick={() => startTransition(() => submit(channelURL))}
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

export default TwitchDialog