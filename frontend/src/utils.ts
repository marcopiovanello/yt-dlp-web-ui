import { blue, red } from '@mui/material/colors'
import { pipe } from 'fp-ts/lib/function'
import { Accent, ThemeNarrowed } from './atoms/settings'
import type { RPCResponse } from "./types"
import { ProcessStatus } from './types'

export function validateIP(ipAddr: string): boolean {
  let ipRegex = /^(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]\d|\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]\d|\d)){3}$/gm
  return ipRegex.test(ipAddr)
}

export function validateDomain(url: string): boolean {
  const urlRegex = /(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*)/
  const slugRegex = /^[a-z0-9]+(?:-[a-z0-9]+)*$/

  const [name, slug] = url.split('/')

  return urlRegex.test(url) || name === 'localhost' && slugRegex.test(slug)
}

export const ellipsis = (str: string, lim: number) =>
  str.length > lim
    ? `${str.substring(0, lim)}...`
    : str

export function toFormatArgs(codes: string[]): string {
  if (codes.length > 1) {
    return codes.reduce((v, a) => ` -f ${v}+${a}`)
  }
  if (codes.length === 1) {
    return ` -f ${codes[0]}`
  }
  return ''
}

export function formatSize(bytes: number): string {
  const threshold = 1024
  const units = ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

  let i = 0
  while (bytes >= threshold) {
    bytes /= threshold
    i = i + 1
  }

  return `${bytes.toFixed(i == 0 ? 0 : 2)} ${units.at(i)}`
}

export const formatSpeedMiB = (val: number) =>
  `${(val / 1_048_576).toFixed(2)} MiB/s`

export const datetimeCompareFunc = (a: string, b: string) =>
  new Date(a).getTime() - new Date(b).getTime()

export function isRPCResponse(object: any): object is RPCResponse<any> {
  return 'result' in object && 'id' in object
}

export function mapProcessStatus(status: ProcessStatus) {
  switch (status) {
    case ProcessStatus.PENDING:
      return 'Pending'
    case ProcessStatus.DOWNLOADING:
      return 'Downloading'
    case ProcessStatus.COMPLETED:
      return 'Completed'
    case ProcessStatus.ERRORED:
      return 'Error'
    case ProcessStatus.LIVESTREAM:
      return 'Livestream'
    default:
      return 'Pending'
  }
}

export const prefersDarkMode = () =>
  window.matchMedia('(prefers-color-scheme: dark)').matches

export const base64URLEncode = (s: string) => pipe(
  s,
  s => String.fromCodePoint(...new TextEncoder().encode(s)),
  btoa,
  encodeURIComponent
)

export const base64URLDecode = (encoded: string): string => {
  try {
    const decoded = decodeURIComponent(encoded)
    const bytes = Uint8Array.from(atob(decoded), c => c.charCodeAt(0))
    return new TextDecoder().decode(bytes)
  } catch {
    return ''
  }
}

export const getAccentValue = (accent: Accent, mode: ThemeNarrowed) => {
  switch (accent) {
    case 'default':
      return mode === 'light' ? blue[700] : blue[300]
    case 'red':
      return mode === 'light' ? red[600] : red[400]
    default:
      return mode === 'light' ? blue[700] : blue[300]
  }
}