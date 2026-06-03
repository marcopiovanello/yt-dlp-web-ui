import { pipe } from 'fp-ts/lib/function'
import { of } from 'fp-ts/lib/Task'
import { getOrElse } from 'fp-ts/lib/TaskEither'
import { ffetch } from '../lib/httpClient'
import { RPCVersion } from '../types'
import { rpcClientState } from './rpc'
import { serverURL } from './settings'
import { atom } from 'jotai'

export const connectedState = atom(false)

export const freeSpaceBytesState = atom(async (get) => {
  const res = await get(rpcClientState)
    .freeSpace()
    .catch(() => ({ result: 0 }))
  return res.result
})

export const availableDownloadPathsState = atom(async (get) => {
  const res = await get(rpcClientState).directoryTree()
    .catch(() => ({ result: [] }))
  return res.result
})

export const ytdlpRpcVersionState = atom<Promise<RPCVersion>>(async (get) => await pipe(
  ffetch<RPCVersion>(`${get(serverURL)}/api/v1/version`),
  getOrElse(() => pipe(
    {
      rpcVersion: 'unknown version',
      ytdlpVersion: 'unknown version',
    },
    of
  )),
)())