import { pipe } from 'fp-ts/lib/function'
import { matchW } from 'fp-ts/lib/TaskEither'
import { ffetch } from '../lib/httpClient'
import { prefersDarkMode } from '../utils'
import { atomWithStorage } from 'jotai/utils'
import { atom } from 'jotai'

export const languages = [
  'catalan',
  'chinese',
  'english',
  'french',
  'german',
  'italian',
  'japanese',
  'korean',
  'polish',
  'portuguese-br',
  'russian',
  'spanish',
  'swedish',
  'ukrainian'
] as const

export type Language = (typeof languages)[number]

export type Theme = 'light' | 'dark' | 'system'
export type ThemeNarrowed = 'light' | 'dark'

export const accents = ['default', 'red'] as const
export type Accent = (typeof accents)[number]

export interface SettingsState {
  serverAddr: string
  serverPort: number
  language: Language
  theme: ThemeNarrowed
  accent: Accent
  cliArgs: string
  formatSelection: boolean
  fileRenaming: boolean
  pathOverriding: boolean
  enableCustomArgs: boolean
  listView: boolean
  servedFromReverseProxy: boolean
  appTitle: string
}

export const languageState = atomWithStorage<Language>(
  'language',
  localStorage.getItem('language') as Language || 'english'
)

export const themeState = atomWithStorage<Theme>(
  'theme',
  localStorage.getItem('theme') as Theme || 'system'
)

export const serverAddressState = atomWithStorage<string>(
  'server-addr',
  localStorage.getItem('server-addr') || window.location.hostname
)

export const serverPortState = atomWithStorage<number>(
  'server-port',
  Number(localStorage.getItem('server-port')) || Number(window.location.port)
)

export const latestCliArgumentsState = atomWithStorage<string>(
  'cli-args',
  localStorage.getItem('cli-args') || '--no-mtime'
)

export const formatSelectionState = atomWithStorage(
  'format-selection',
  localStorage.getItem('format-selection') === 'true'
)

export const fileRenamingState = atomWithStorage(
  'file-renaming',
  localStorage.getItem('file-renaming') === 'true'
)

export const pathOverridingState = atomWithStorage(
  'path-overriding',
  localStorage.getItem('path-overriding') === 'true'
)

export const enableCustomArgsState = atomWithStorage(
  'enable-custom-args',
  localStorage.getItem('enable-custom-args') === 'true'
)

export const listViewState = atomWithStorage(
  'listview',
  localStorage.getItem('listview') === 'true'
)

export const servedFromReverseProxyState = atomWithStorage(
  'reverseProxy',
  localStorage.getItem('reverseProxy') === 'true' || window.location.port == ''
)

export const servedFromReverseProxySubDirState = atomWithStorage<string>(
  'reverseProxySubDir',
  localStorage.getItem('reverseProxySubDir') ?? ''
)

export const appTitleState = atomWithStorage(
  'appTitle',
  localStorage.getItem('appTitle') ?? 'yt-dlp Web UI'
)

export const serverAddressAndPortState = atom((get) => {
  if (get(servedFromReverseProxySubDirState)) {
    return `${get(serverAddressState)}/${get(servedFromReverseProxySubDirState)}/`
  }
  if (get(servedFromReverseProxyState)) {
    return `${get(serverAddressState)}`
  }
  return `${get(serverAddressState)}:${get(serverPortState)}`
})

export const serverURL = atom((get) =>
  `${window.location.protocol}//${get(serverAddressAndPortState)}`
)

export const rpcWebSocketEndpoint = atom((get) => {
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${proto}//${get(serverAddressAndPortState)}/rpc/ws`
}
)

export const rpcHTTPEndpoint = atom((get) => {
  const proto = window.location.protocol
  return `${proto}//${get(serverAddressAndPortState)}/rpc/http`
}
)

export const serverSideCookiesState = atom<Promise<string>>(async (get) => await pipe(
  ffetch<Readonly<{ cookies: string }>>(`${get(serverURL)}/api/v1/cookies`),
  matchW(
    () => '',
    (r) => r.cookies
  )
)())

const themeSelector = atom<ThemeNarrowed>((get) => {
  const theme = get(themeState)
  if ((theme === 'system' && prefersDarkMode()) || theme === 'dark') {
    return 'dark'
  }
  return 'light'
})

export const accentState = atomWithStorage<Accent>(
  'accent-color',
  localStorage.getItem('accent-color') as Accent ?? 'default',
)

export const settingsState = atom<SettingsState>((get) => ({
  serverAddr: get(serverAddressState),
  serverPort: get(serverPortState),
  language: get(languageState),
  theme: get(themeSelector),
  accent: get(accentState),
  cliArgs: get(latestCliArgumentsState),
  formatSelection: get(formatSelectionState),
  fileRenaming: get(fileRenamingState),
  pathOverriding: get(pathOverridingState),
  enableCustomArgs: get(enableCustomArgsState),
  listView: get(listViewState),
  servedFromReverseProxy: get(servedFromReverseProxyState),
  appTitle: get(appTitleState)
})
)
