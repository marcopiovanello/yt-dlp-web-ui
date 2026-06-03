import * as O from 'fp-ts/Option'
import { useEffect, useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import { take, timer } from 'rxjs'
import { downloadsState } from '../atoms/downloads'
import { rpcPollingTimeState } from '../atoms/rpc'
import { serverAddressAndPortState } from '../atoms/settings'
import { connectedState } from '../atoms/status'
import { useSubscription } from '../hooks/observable'
import { useToast } from '../hooks/toast'
import { useI18n } from '../hooks/useI18n'
import { useRPC } from '../hooks/useRPC'
import { datetimeCompareFunc, isRPCResponse } from '../utils'
import { useAtom, useAtomValue, useSetAtom } from 'jotai'

interface Props extends React.HTMLAttributes<HTMLBaseElement> { }

const SocketSubscriber: React.FC<Props> = () => {
  const [connected, setIsConnected] = useAtom(connectedState)
  const setDownloads = useSetAtom(downloadsState)

  const serverAddressAndPort = useAtomValue(serverAddressAndPortState)
  const rpcPollingTime = useAtomValue(rpcPollingTimeState)

  const { i18n } = useI18n()
  const { client } = useRPC()
  const { pushMessage } = useToast()

  const navigate = useNavigate()

  const socketOnce$ = useMemo(() => client.socket$.pipe(take(1)), [])

  useEffect(() => {
    if (!connected) {
      socketOnce$.subscribe(() => {
        setIsConnected(true)
        pushMessage(
          `${i18n.t('toastConnected')} (${serverAddressAndPort})`,
          "success"
        )
      })
    }
  }, [connected])

  useSubscription(
    client.socket$,
    event => {
      if (!isRPCResponse(event)) { return }
      if (!Array.isArray(event.result)) { return }

      if (event.result) {
        return setDownloads(
          O.of(event.result
            .filter(f => !!f.info.url).sort((a, b) => datetimeCompareFunc(
              b.info.created_at,
              a.info.created_at,
            )),
          )
        )
      }
      setDownloads(O.none)
    },
    err => {
      console.error(err)
      pushMessage(
        `${i18n.t('rpcConnErr')} (${serverAddressAndPort})`,
        "error"
      ),
        navigate(`/error`)
    }
  )

  useEffect(() => {
    if (connected) {
      const sub = timer(0, rpcPollingTime).subscribe(() => client.running())
      return () => sub.unsubscribe()
    }
  }, [connected, client, rpcPollingTime])

  return null
}

export default SocketSubscriber