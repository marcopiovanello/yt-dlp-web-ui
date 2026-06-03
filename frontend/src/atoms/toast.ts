import { AlertColor } from '@mui/material'
import { atom } from 'jotai'

export type Toast = {
  open: boolean,
  message: string
  autoClose: boolean
  createdAt: number,
  severity?: AlertColor
}

export const toastListState = atom<Toast[]>([])