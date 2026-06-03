import AddIcon from '@mui/icons-material/Add'
import CloseIcon from '@mui/icons-material/Close'
import {
  Alert,
  AppBar,
  Backdrop,
  Box,
  Button,
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
import { forwardRef, useEffect, useState, useTransition } from 'react'
import { serverURL } from '../atoms/settings'
import { useToast } from '../hooks/toast'
import { useI18n } from '../hooks/useI18n'
import { ffetch } from '../lib/httpClient'
import { CustomTemplate } from '../types'
import TemplateTextField from './TemplateTextField'

const Transition = forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement
  },
  ref: React.Ref<unknown>,
) {
  return <Slide direction="up" ref={ref} {...props} />
})

interface Props extends React.HTMLAttributes<HTMLBaseElement> {
  open: boolean
  onClose: () => void
}

const TemplatesEditor: React.FC<Props> = ({ open, onClose }) => {
  const [templateName, setTemplateName] = useState('')
  const [templateContent, setTemplateContent] = useState('')

  const serverAddr = useAtomValue(serverURL)
  const [isPending, startTransition] = useTransition()

  const [templates, setTemplates] = useState<CustomTemplate[]>([])

  const { i18n } = useI18n()
  const { pushMessage } = useToast()

  useEffect(() => {
    if (open) {
      fetchTemplates()
    }
  }, [open])

  const fetchTemplates = async () => {
    const task = ffetch<CustomTemplate[]>(`${serverAddr}/api/v1/template/all`)
    const either = await task()

    pipe(
      either,
      matchW(
        (l) => pushMessage(l),
        (r) => setTemplates(r)
      )
    )
  }

  const addTemplate = async () => {
    const task = ffetch<unknown>(`${serverAddr}/api/v1/template`, {
      method: 'POST',
      body: JSON.stringify({
        name: templateName,
        content: templateContent,
      })
    })

    const either = await task()

    pipe(
      either,
      matchW(
        (l) => pushMessage(l, 'warning'),
        () => {
          pushMessage('Added template')
          fetchTemplates()
          setTemplateName('')
          setTemplateContent('')
        }
      )
    )
  }

  const updateTemplate = async (template: CustomTemplate) => {
    const task = ffetch<CustomTemplate>(`${serverAddr}/api/v1/template`, {
      method: 'PATCH',
      body: JSON.stringify(template)
    })

    const either = await task()

    pipe(
      either,
      matchW(
        (l) => pushMessage(l, 'warning'),
        (r) => {
          pushMessage(`Updated template ${r.name}`)
          fetchTemplates()
        }
      )
    )
  }

  const deleteTemplate = async (id: string) => {
    const task = ffetch<unknown>(`${serverAddr}/api/v1/template/${id}`, {
      method: 'DELETE',
    })

    const either = await task()

    pipe(
      either,
      matchW(
        (l) => pushMessage(l, 'warning'),
        () => {
          pushMessage('Deleted template')
          fetchTemplates()
        }
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
      <Backdrop
        sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
        open={isPending}
      />
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
            {i18n.t('templatesEditor')}
          </Typography>
        </Toolbar>
      </AppBar>
      <Box sx={{
        backgroundColor: (theme) => theme.palette.background.default,
        minHeight: (theme) => `calc(99vh - ${theme.mixins.toolbar.minHeight}px)`
      }}>
        <Grid container spacing={2} sx={{ px: 4, pt: 3, pb: 4 }}>
          <Grid item>
            <Alert severity="info">
              {i18n.t('templatesReloadInfo')}
            </Alert>
          </Grid>
          <Grid item xs={12}>
            <Paper
              elevation={4}
              sx={{
                p: 2,
                display: 'flex',
                flexDirection: 'column',
              }}
            >
              <Grid container spacing={2} justifyContent="center" alignItems="center">
                <Grid item xs={3}>
                  <TextField
                    fullWidth
                    label={i18n.t('templatesEditorNameLabel')}
                    onChange={e => setTemplateName(e.currentTarget.value)}
                    value={templateName}
                  />
                </Grid>
                <Grid item xs={9}>
                  <TextField
                    fullWidth
                    label={i18n.t('templatesEditorContentLabel')}
                    onChange={e => setTemplateContent(e.currentTarget.value)}
                    value={templateContent}
                    InputProps={{
                      endAdornment: <Button
                        variant='contained'
                        onClick={() => startTransition(async () => await addTemplate())}
                      >
                        <AddIcon />
                      </Button>
                    }}
                  />
                </Grid>
              </Grid>
              {templates.map(template => (
                <TemplateTextField
                  key={template.id}
                  template={template}
                  onChange={updateTemplate}
                  onDelete={deleteTemplate}
                />
              ))}
            </Paper>
          </Grid>
        </Grid>
      </Box>
    </Dialog >
  )
}

export default TemplatesEditor