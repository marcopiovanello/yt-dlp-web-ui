import { Container, Paper, Typography } from '@mui/material'
import LogTerminal from '../components/LogTerminal'
import { useI18n } from '../hooks/useI18n'

const Terminal: React.FC = () => {
  const { i18n } = useI18n()

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      <Paper
        sx={{
          p: 2.5,
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <Typography pb={2} variant="h5" color="primary">
          {i18n.t('logsTitle')}
        </Typography>
        <LogTerminal />
      </Paper>
    </Container >
  )
}

export default Terminal