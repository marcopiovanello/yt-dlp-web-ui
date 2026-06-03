import { atom } from 'jotai'
import { RPCClient } from '../lib/rpcClient'
import { rpcHTTPEndpoint, rpcWebSocketEndpoint } from './settings'
import { atomWithStorage } from 'jotai/utils'

export const rpcClientState = atom((get) =>
  new RPCClient(
    get(rpcHTTPEndpoint),
    get(rpcWebSocketEndpoint),
    localStorage.getItem('token') ?? ''
  ),
)

export const rpcPollingTimeState = atomWithStorage(
  'rpc-polling-time',
  Number(localStorage.getItem('rpc-polling-time')) || 1000
)