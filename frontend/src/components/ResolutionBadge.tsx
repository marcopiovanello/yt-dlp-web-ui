import EightK from '@mui/icons-material/EightK'
import FourK from '@mui/icons-material/FourK'
import Hd from '@mui/icons-material/Hd'
import Sd from '@mui/icons-material/Sd'

const ResolutionBadge: React.FC<{ resolution?: string }> = ({ resolution }) => {
  if (!resolution) return null
  if (resolution.includes('4320')) return <EightK color="primary" />
  if (resolution.includes('2160')) return <FourK color="primary" />
  if (resolution.includes('1080')) return <Hd color="primary" />
  if (resolution.includes('720')) return <Sd color="primary" />
  return null
}

export default ResolutionBadge