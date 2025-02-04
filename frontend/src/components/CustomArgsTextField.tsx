import { TextField } from '@mui/material'
import { useAtom, useAtomValue } from 'jotai'
import { customArgsState } from '../atoms/downloadTemplate'
import { settingsState } from '../atoms/settings'
import { useI18n } from '../hooks/useI18n'
import { useEffect } from 'react'

const CustomArgsTextField: React.FC = () => {
  const { i18n } = useI18n()

  const settings = useAtomValue(settingsState)

  const [customArgs, setCustomArgs] = useAtom(customArgsState)

  useEffect(() => {
    setCustomArgs('')
  }, [])

  const handleCustomArgsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCustomArgs(e.target.value)
  }

  return (
    <TextField
      fullWidth
      label={i18n.t('customArgsInput')}
      variant="outlined"
      onChange={handleCustomArgsChange}
      value={customArgs}
      disabled={settings.formatSelection}
    />
  )
}

export default CustomArgsTextField