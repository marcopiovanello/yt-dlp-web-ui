import { ThemeProvider } from '@emotion/react'
import ChevronLeft from '@mui/icons-material/ChevronLeft'
import CloudDownloadIcon from '@mui/icons-material/CloudDownload'
import Dashboard from '@mui/icons-material/Dashboard'
import LiveTvIcon from '@mui/icons-material/LiveTv'
import Menu from '@mui/icons-material/Menu'
import SettingsIcon from '@mui/icons-material/Settings'
import TerminalIcon from '@mui/icons-material/Terminal'
import UpdateIcon from '@mui/icons-material/Update'
import { Box, createTheme } from '@mui/material'
import CssBaseline from '@mui/material/CssBaseline'
import Divider from '@mui/material/Divider'
import IconButton from '@mui/material/IconButton'
import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemIcon from '@mui/material/ListItemIcon'
import ListItemText from '@mui/material/ListItemText'
import Toolbar from '@mui/material/Toolbar'
import Typography from '@mui/material/Typography'
import { grey } from '@mui/material/colors'
import { useAtomValue } from 'jotai'
import { useMemo, useState } from 'react'
import { Link, Outlet, useLocation } from 'react-router-dom'
import { settingsState } from './atoms/settings'
import AppBar from './components/AppBar'
import Drawer from './components/Drawer'
import Footer from './components/Footer'
import Logout from './components/Logout'
import SocketSubscriber from './components/SocketSubscriber'
import ThemeToggler from './components/ThemeToggler'
import { useI18n } from './hooks/useI18n'
import Toaster from './providers/ToasterProvider'
import { getAccentValue } from './utils'

export default function Layout() {
  const [open, setOpen] = useState(false)
  const location = useLocation()
  const isPlayerRoute = location.pathname.startsWith('/public/')

  const settings = useAtomValue(settingsState)
  const mode = settings.theme

  const theme = useMemo(() =>
    createTheme({
      palette: {
        mode: settings.theme,
        primary: {
          main: getAccentValue(settings.accent, settings.theme)
        },
        background: {
          default: settings.theme === 'light' ? grey[50] : '#121212'
        },
      },
    }), [settings.theme, settings.accent]
  )

  const toggleDrawer = () => setOpen(state => !state)
  const { i18n } = useI18n()

  return (
    <ThemeProvider theme={theme}>
      <title>{settings.appTitle}</title>
      <SocketSubscriber />
      <CssBaseline />

      {isPlayerRoute ? (
        // ❯ Minimal layout for /public/:encoded
        <Box sx={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
          <Outlet />
        </Box>
      ) : (
        // ❯ Full layout for all other pages
        <Box sx={{ display: 'flex' }}>
          <AppBar position="absolute" open={open}>
            <Toolbar sx={{ pr: '24px' }}>
              <IconButton
                edge="start"
                color="inherit"
                aria-label="open drawer"
                onClick={toggleDrawer}
                sx={{
                  marginRight: '36px',
                  ...(open && { display: 'none' }),
                }}
              >
                <Menu />
              </IconButton>
              <Typography
                component="h1"
                variant="h6"
                color="inherit"
                noWrap
                sx={{ flexGrow: 1 }}
              >
                {settings.appTitle}
              </Typography>
            </Toolbar>
          </AppBar>
          <Drawer variant="permanent" open={open}>
            <Toolbar
              sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'flex-end',
                px: [1],
              }}
            >
              <IconButton onClick={toggleDrawer}>
                <ChevronLeft />
              </IconButton>
            </Toolbar>
            <Divider />
            <List component="nav">
              <Link to={'/'} style={{ textDecoration: 'none', color: mode === 'dark' ? '#ffffff' : '#000000DE' }}>
                <ListItemButton>
                  <ListItemIcon><Dashboard /></ListItemIcon>
                  <ListItemText primary={i18n.t('homeButtonLabel')} />
                </ListItemButton>
              </Link>
              <Link to={'/filebrowser'} style={{ textDecoration: 'none', color: mode === 'dark' ? '#ffffff' : '#000000DE' }}>
                <ListItemButton>
                  <ListItemIcon><CloudDownloadIcon /></ListItemIcon>
                  <ListItemText primary={i18n.t('archiveButtonLabel')} />
                </ListItemButton>
              </Link>
              <Link to={'/subscriptions'} style={{ textDecoration: 'none', color: mode === 'dark' ? '#ffffff' : '#000000DE' }}>
                <ListItemButton>
                  <ListItemIcon><UpdateIcon /></ListItemIcon>
                  <ListItemText primary={i18n.t('subscriptionsButtonLabel')} />
                </ListItemButton>
              </Link>
              <Link to={'/monitor'} style={{ textDecoration: 'none', color: mode === 'dark' ? '#ffffff' : '#000000DE' }}>
                <ListItemButton>
                  <ListItemIcon><LiveTvIcon /></ListItemIcon>
                  <ListItemText primary={i18n.t('archiveButtonLabel')} />
                </ListItemButton>
              </Link>
              <Link to={'/log'} style={{ textDecoration: 'none', color: mode === 'dark' ? '#ffffff' : '#000000DE' }}>
                <ListItemButton>
                  <ListItemIcon><TerminalIcon /></ListItemIcon>
                  <ListItemText primary={i18n.t('logsTitle')} />
                </ListItemButton>
              </Link>
              <Link to={'/settings'} style={{ textDecoration: 'none', color: mode === 'dark' ? '#ffffff' : '#000000DE' }}>
                <ListItemButton>
                  <ListItemIcon><SettingsIcon /></ListItemIcon>
                  <ListItemText primary={i18n.t('settingsButtonLabel')} />
                </ListItemButton>
              </Link>
              <ThemeToggler />
              <Logout />
            </List>
          </Drawer>
          <Box
            component="main"
            sx={{
              flexGrow: 1,
              height: '100vh',
              overflow: 'auto',
            }}
          >
            <Toolbar />
            <Outlet />
          </Box>
        </Box>
      )}

      {!isPlayerRoute && <Footer />}
      <Toaster />
    </ThemeProvider>
  )
}
