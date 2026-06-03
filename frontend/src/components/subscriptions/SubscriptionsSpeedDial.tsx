import AddCircleIcon from '@mui/icons-material/AddCircle'
import { SpeedDial, SpeedDialAction, SpeedDialIcon } from '@mui/material'
import { useI18n } from '../../hooks/useI18n'

type Props = {
  onOpen: () => void
}

const SubscriptionsSpeedDial: React.FC<Props> = ({ onOpen }) => {
  const { i18n } = useI18n()

  return (
    <SpeedDial
      ariaLabel="Subscriptions speed dial"
      sx={{ position: 'absolute', bottom: 64, right: 24 }}
      icon={<SpeedDialIcon />}
    >
      <SpeedDialAction
        icon={<AddCircleIcon />}
        tooltipTitle={i18n.t('newSubscriptionButton')}
        onClick={onOpen}
      />
    </SpeedDial>
  )
}

export default SubscriptionsSpeedDial