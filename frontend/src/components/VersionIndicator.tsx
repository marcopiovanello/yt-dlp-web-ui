import { Chip } from '@mui/material'
import { ytdlpRpcVersionState } from '../atoms/status'
import { useAtomValue } from 'jotai'

const VersionIndicator: React.FC = () => {
  const version = useAtomValue(ytdlpRpcVersionState)

  return (
    <div style={{ display: 'flex', gap: 4, alignItems: 'center' }}>
      <Chip label={`UI v3.2.5`} variant="outlined" size="small" />
      <Chip label={`RPC v${version.rpcVersion}`} variant="outlined" size="small" />
      <Chip label={`yt-dlp v${version.ytdlpVersion}`} variant="outlined" size="small" />
    </div>
  )
}

export default VersionIndicator