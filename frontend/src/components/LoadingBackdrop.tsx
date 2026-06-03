import { Backdrop, CircularProgress } from '@mui/material'

type Props = {
  isLoading: boolean
}

const LoadingBackdrop: React.FC<Props> = ({ isLoading }) => {
  return (
    <Backdrop
      sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
      open={isLoading}
    >
      <CircularProgress color="primary" />
    </Backdrop>
  )
}

export default LoadingBackdrop