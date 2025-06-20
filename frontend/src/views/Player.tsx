import { useParams } from 'react-router-dom'
import { useAtomValue } from 'jotai'
import { useEffect, useState } from 'react'
import { serverURL } from '../atoms/settings'
import { base64URLDecode } from '../utils'

import {
  Box,
  CircularProgress,
  Alert
} from '@mui/material'

const Player: React.FC = () => {
  const { encoded } = useParams<{ encoded: string }>()
  const serverAddr = useAtomValue(serverURL)

  const [videoSrc, setVideoSrc] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!encoded) {
      setError('No video specified.')
      setLoading(false)
      return
    }

    try {
      base64URLDecode(encoded) 
      const url = `${serverAddr}/public/${encoded}`
      setVideoSrc(url)
    } catch (e) {
      setError('Invalid video path.')
    } finally {
      setLoading(false)
    }
  }, [encoded, serverAddr])

  return (
    <Box
      sx={{
        width: '100vw',
        height: '100dvh', 
        bgcolor: 'black',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        overflow: 'hidden'
      }}
    >
      {loading && <CircularProgress color="inherit" />}
      {error && <Alert severity="error">{error}</Alert>}
      {videoSrc && (
        <video
          controls
          autoPlay
          src={videoSrc}
          style={{
            width: '100%',
            height: '100%',
            objectFit: 'contain',
            backgroundColor: 'black',
          }}
        />
      )}
    </Box>
  )
}

export default Player
