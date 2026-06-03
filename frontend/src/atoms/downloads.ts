import * as O from 'fp-ts/Option'
import { pipe } from 'fp-ts/lib/function'
import { RPCResult } from '../types'
import { atom } from 'jotai'

export const downloadsState = atom<O.Option<RPCResult[]>>(O.none)

export const loadingDownloadsState = atom<boolean>((get) => O.isNone(get(downloadsState)))

export const activeDownloadsState = atom<RPCResult[]>((get) => pipe(
  get(downloadsState),
  O.getOrElse(() => new Array<RPCResult>())
))