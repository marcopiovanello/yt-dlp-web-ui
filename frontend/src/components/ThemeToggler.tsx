import Brightness4 from '@mui/icons-material/Brightness4'
import Brightness5 from '@mui/icons-material/Brightness5'
import BrightnessAuto from '@mui/icons-material/BrightnessAuto'
import { ListItemButton, ListItemIcon, ListItemText } from '@mui/material'
import { Theme, themeState } from '../atoms/settings'
import { useI18n } from '../hooks/useI18n'
import { useAtom } from 'jotai'

const ThemeToggler: React.FC = () => {
  const [theme, setTheme] = useAtom(themeState)

  const actions: Record<Theme, React.ReactNode> = {
    system: <BrightnessAuto />,
    light: <Brightness4 />,
    dark: <Brightness5 />,
  }

  const themes: Theme[] = ['system', 'light', 'dark']
  const currentTheme = themes.indexOf(theme)

  const { i18n } = useI18n()

  return (
    <ListItemButton onClick={() => {
      setTheme(themes[(currentTheme + 1) % themes.length])
    }}>
      <ListItemIcon>
        {actions[theme]}
      </ListItemIcon>
      <ListItemText primary={i18n.t('themeTogglerLabel')} />
    </ListItemButton>
  )
}

export default ThemeToggler