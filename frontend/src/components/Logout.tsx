import LogoutIcon from '@mui/icons-material/Logout'
import { ListItemButton, ListItemIcon, ListItemText } from '@mui/material'
import { useNavigate } from 'react-router-dom'
import { useI18n } from '../hooks/useI18n'

export default function Logout() {
  const navigate = useNavigate()

  const logout = async () => {
    localStorage.removeItem('token')
    navigate('/login')
  }

  const { i18n } = useI18n()

  return (
    <ListItemButton onClick={logout}>
      <ListItemIcon>
        <LogoutIcon />
      </ListItemIcon>
      <ListItemText primary={i18n.t('rpcAuthenticationLabel')} />
    </ListItemButton>
  )
}