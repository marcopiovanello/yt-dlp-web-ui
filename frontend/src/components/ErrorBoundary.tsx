import ErrorIcon from '@mui/icons-material/Error'
import { Button, Container, SvgIcon, Typography, styled } from '@mui/material'
import { Link } from 'react-router-dom'

const FlexContainer = styled(Container)({
  display: 'flex',
  minWidth: '100%',
  minHeight: '80vh',
  alignItems: 'center',
  justifyContent: 'center',
  flexDirection: 'column'
})

const Title = styled(Typography)({
  display: 'flex',
  width: '100%',
  alignItems: 'center',
  justifyContent: 'center',
  paddingBottom: '0.5rem'
})

const ErrorBoundary: React.FC = () => {
  return (
    <FlexContainer>
      <Title fontWeight={'500'} fontSize={72} color={'gray'}>
        <SvgIcon sx={{ fontSize: '200px' }}>
          <ErrorIcon />
        </SvgIcon>
      </Title>
      <Title fontWeight={'500'} fontSize={36} color={'gray'}>
        An error occurred :\
      </Title>
      <Title fontWeight={'400'} fontSize={28} color={'gray'}>
        Check your settings!
      </Title>
      <Link to={'/settings'} >
        <Button variant='contained' sx={{ mt: 2 }}>
          Goto Settings
        </Button>
      </Link>
      <Typography sx={{ mt: 2 }} color={'gray'} fontWeight={500}>
        Or login if authentification is enabled
      </Typography>
      <Link to={'/login'} >
        <Button variant='contained' sx={{ mt: 2 }}>
          login
        </Button>
      </Link>
    </FlexContainer>
  )
}

export default ErrorBoundary
