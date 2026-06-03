import { atom } from 'jotai'
import { activeDownloadsState } from './downloads'

export const loadingAtom = atom(true)

export const totalDownloadSpeedState = atom<number>((get) =>
  get(activeDownloadsState)
    .map(d => d.progress.speed)
    .reduce((curr, next) => curr + next, 0)
)