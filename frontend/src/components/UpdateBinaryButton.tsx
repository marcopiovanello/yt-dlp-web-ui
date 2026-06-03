import { Button, CircularProgress } from '@mui/material'
import { useI18n } from '../hooks/useI18n'
import { useRPC } from '../hooks/useRPC'
import { useState } from 'react'
import { useToast } from '../hooks/toast'

const UpdateBinaryButton: React.FC = () => {
  const { i18n } = useI18n()
  const { client } = useRPC()
  const { pushMessage } = useToast()

  const [isLoading, setIsLoading] = useState(false)

  const updateBinary = () => {
    setIsLoading(true)
    client
      .updateExecutable()
      .then(() => pushMessage(i18n.t('toastUpdated'), 'success'))
      .then(() => setIsLoading(false))
  }

  return (
    <Button
      variant="contained"
      endIcon={isLoading ? <CircularProgress size={16} color='secondary' /> : <></>}
      onClick={updateBinary}
    >
      {i18n.t('updateBinButton')}
    </Button>
  )
}

export default UpdateBinaryButton