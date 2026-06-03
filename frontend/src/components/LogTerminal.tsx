import { useEffect, useMemo, useRef, useState } from 'react'
import { serverURL } from '../atoms/settings'
import { useI18n } from '../hooks/useI18n'
import { useAtomValue } from 'jotai'

const token = localStorage.getItem('token')

const LogTerminal: React.FC = () => {
  const [logBuffer, setLogBuffer] = useState<string[]>([])
  const [isConnecting, setIsConnecting] = useState(true)

  const boxRef = useRef<HTMLDivElement>(null)

  const serverAddr = useAtomValue(serverURL)

  const { i18n } = useI18n()

  const eventSource = useMemo(
    () => new EventSource(`${serverAddr}/log/sse?token=${token}`),
    [serverAddr]
  )

  useEffect(() => {
    eventSource.addEventListener('log', event => {
      const msg: string = JSON.parse(event.data)
      setLogBuffer(buff => [...buff, msg].slice(-500))

      boxRef.current?.scrollTo(0, boxRef.current.scrollHeight)
    })

    // TODO: in dev mode it breaks sse
    return () => eventSource.close()
  }, [eventSource])

  useEffect(() => {
    eventSource.onopen = () => setIsConnecting(false)
  }, [eventSource])

  const logEntryStyle = (data: string) => {
    const sx = {}

    if (data.includes("level=ERROR")) {
      return { ...sx, color: 'red' }
    }
    if (data.includes("level=WARN")) {
      return { ...sx, color: 'orange' }
    }

    return sx
  }

  return (

    <div
      ref={boxRef}
      style={{
        fontFamily: 'Roboto Mono',
        height: '70.5vh',
        overflowY: 'auto',
        overflowX: 'auto',
        fontSize: '13.5px',
        fontWeight: '600',
        backgroundColor: 'black',
        color: 'white',
        padding: '0.5rem',
        borderRadius: '0.25rem'
      }}
    >
      {isConnecting ? <div>{'Connecting...'}</div> : <div>{'Connected!'}</div>}

      {logBuffer.length === 0 && <div>{i18n.t('awaitingLogs')}</div>}

      {logBuffer.map((log, idx) => (
        <div key={idx} style={logEntryStyle(log)}>
          {log}
        </div>
      ))}
    </div>

  )
}

export default LogTerminal