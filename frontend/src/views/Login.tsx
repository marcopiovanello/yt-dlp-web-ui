/*
  Login view component
*/

import styled from '@emotion/styled'
import {
  Button,
  Container,
  Divider,
  Paper,
  Stack,
  TextField,
  Typography
} from '@mui/material'
import { matchW } from 'fp-ts/lib/TaskEither'
import { pipe } from 'fp-ts/lib/function'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { serverURL } from '../atoms/settings'
import { useToast } from '../hooks/toast'
import { ffetch } from '../lib/httpClient'
import { useAtomValue } from 'jotai'

const LoginContainer = styled(Container)({
  display: 'flex',
  minWidth: '100%',
  minHeight: '85vh',
  alignItems: 'center',
  justifyContent: 'center',
})

const Title = styled(Typography)({
  display: 'flex',
  width: '100%',
  alignItems: 'center',
  justifyContent: 'center',
  paddingBottom: '0.5rem'
})

export default function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')

  const [formHasError, setFormHasError] = useState(false)

  const url = useAtomValue(serverURL)

  const navigate = useNavigate()

  const { pushMessage } = useToast()

  const navigateAndReload = () => {
    navigate('/')
    window.location.reload()
  }

  const login = async () => {
    const task = ffetch<string>(`${url}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username,
        password,
      }),
    })

    pipe(
      task,
      matchW(
        (error) => {
          setFormHasError(true)
          pushMessage(error, 'error')
        },
        (token) => {
          console.log(token)
          localStorage.setItem('token', token)
          navigateAndReload()
        }
      )
    )()
  }

  const loginWithOpenId = () => window.open(`${url}/auth/openid/login`)

  return (
    <LoginContainer>
      <Paper sx={{ padding: '1.5rem', minWidth: '25%' }}>
        <Stack direction="column" spacing={2}>
          <Title fontWeight={'700'} fontSize={32} color={'primary'}>
            yt-dlp WebUI
          </Title>
          <Title fontWeight={'500'} fontSize={16} color={'gray'}>
            To configure authentication check the&nbsp;
            <a href='https://github.com/marcopiovanello/yt-dlp-web-ui/wiki/Authentication-methods'>wiki</a>.
          </Title>
          <TextField
            label="Username"
            type="text"
            autoComplete="yt-dlp-webui-username"
            error={formHasError}
            onChange={e => setUsername(e.currentTarget.value)}
          />
          <TextField
            label="Password"
            type="password"
            autoComplete="yt-dlp-webui-password"
            error={formHasError}
            onChange={e => setPassword(e.currentTarget.value)}
          />
          <Button variant="contained" size="large" onClick={() => login()}>
            Submit
          </Button>

          <Divider>
            <Typography color={'gray'}>
              or use your authentication provider
            </Typography>
          </Divider>

          <Button variant="contained" size="large" onClick={loginWithOpenId}>
            Login with OpenID
          </Button>
        </Stack>
      </Paper>
    </LoginContainer>
  )
}