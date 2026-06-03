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
import { forwardRef, useState } from 'react'
import { useToast } from '../../hooks/toast'
import { useI18n } from '../../hooks/useI18n'
import { useRPC } from '../../hooks/useRPC'

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

const LivestreamDialog: React.FC<Props> = ({ open, onClose }) => {
  const [livestreamURL, setLivestreamURL] = useState('')

  const { i18n } = useI18n()
  const { client } = useRPC()
  const { pushMessage } = useToast()

  const exec = (url: string) => client.execLivestream(url)

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
            Livestream monitor
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
                <Grid container>
                  <Grid item xs={12} mb={2}>
                    <Alert severity="info">
                      {i18n.t('livestreamDownloadInfo')}
                    </Alert>
                    <Alert severity="warning" sx={{ mt: 1 }}>
                      {i18n.t('livestreamExperimentalWarning')}
                    </Alert>
                  </Grid>
                  <Grid item xs={12}>
                    <TextField
                      multiline
                      fullWidth
                      label={i18n.t('livestreamURLInput')}
                      variant="outlined"
                      onChange={(e) => setLivestreamURL(e.target.value)}
                    />
                  </Grid>
                  <Grid item>
                    <Button
                      sx={{ mt: 2 }}
                      variant="contained"
                      disabled={livestreamURL === ''}
                      onClick={() => {
                        exec(livestreamURL)
                        onClose()
                        pushMessage(`Monitoring ${livestreamURL}`, 'info')
                      }}
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

export default LivestreamDialog