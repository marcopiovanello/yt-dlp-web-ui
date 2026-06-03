import { AlertColor } from '@mui/material'
import { toastListState } from '../atoms/toast'
import { useSetAtom } from 'jotai'

export const useToast = () => {
  const setToasts = useSetAtom(toastListState)

  return {
    pushMessage: (message: string, severity?: AlertColor) => {
      setToasts(state => [...state, {
        open: true,
        message: message,
        severity: severity,
        autoClose: true,
        createdAt: Date.now()
      }])
    }
  }
}