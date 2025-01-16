import LiveTvIcon from '@mui/icons-material/LiveTv'
import { Container, SvgIcon, Typography, styled } from '@mui/material'
import { useI18n } from '../../hooks/useI18n'

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


export default function NoLivestreams() {
  const { i18n } = useI18n()

  return (
    <FlexContainer>
      <Title fontWeight={'500'} fontSize={72} color={'gray'}>
        <SvgIcon sx={{ fontSize: '200px' }}>
          <LiveTvIcon />
        </SvgIcon>
      </Title>
      <Title fontWeight={'500'} fontSize={36} color={'gray'}>
        {i18n.t('livestreamNoMonitoring')}
      </Title>
    </FlexContainer>
  )
}