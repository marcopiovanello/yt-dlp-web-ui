import { useAtomValue } from 'jotai'
import { rpcClientState } from '../atoms/rpc'

export const useRPC = () => {
  const client = useAtomValue(rpcClientState)

  return {
    client
  }
}