import { CircularProgress } from '@mui/material'
import { Suspense, lazy } from 'react'
import { createHashRouter } from 'react-router-dom'
import Layout from './Layout'
import Terminal from './views/Terminal'

const Home = lazy(() => import('./views/Home'))
const Login = lazy(() => import('./views/Login'))
const Archive = lazy(() => import('./views/Archive'))
const Settings = lazy(() => import('./views/Settings'))
const LiveStream = lazy(() => import('./views/Livestream'))
const Filebrowser = lazy(() => import('./views/Filebrowser'))
const Subscriptions = lazy(() => import('./views/Subscriptions'))

const ErrorBoundary = lazy(() => import('./components/ErrorBoundary'))

export const router = createHashRouter([
  {
    path: '/',
    Component: () => <Layout />,
    children: [
      {
        path: '/',
        element: (
          <Suspense fallback={<CircularProgress />}>
            <Home />
          </Suspense >
        ),
        errorElement: (
          <Suspense fallback={<CircularProgress />}>
            <ErrorBoundary />
          </Suspense >
        )
      },
      {
        path: '/settings',
        element: (
          <Suspense fallback={<CircularProgress />}>
            <Settings />
          </Suspense >
        )
      },
      {
        path: '/log',
        element: (
          <Suspense fallback={<CircularProgress />}>
            <Terminal />
          </Suspense >
        )
      },
      {
        path: '/archive',
        element: (
          <Suspense fallback={<CircularProgress />}>
            <Archive />
          </Suspense >
        ),
        errorElement: (
          <Suspense fallback={<CircularProgress />}>
            <ErrorBoundary />
          </Suspense >
        )
      },
      {
        path: '/filebrowser',
        element: (
          <Suspense fallback={<CircularProgress />}>
            <Filebrowser />
          </Suspense >
        ),
        errorElement: (
          <Suspense fallback={<CircularProgress />}>
            <ErrorBoundary />
          </Suspense >
        )
      },
      {
        path: '/subscriptions',
        element: (
          <Suspense fallback={<CircularProgress />}>
            <Subscriptions />
          </Suspense >
        ),
        errorElement: (
          <Suspense fallback={<CircularProgress />}>
            <ErrorBoundary />
          </Suspense >
        )
      },
      {
        path: '/login',
        element: (
          <Suspense fallback={<CircularProgress />}>
            <Login />
          </Suspense >
        )
      },
      {
        path: '/error',
        element: (
          <Suspense fallback={<CircularProgress />}>
            <ErrorBoundary />
          </Suspense >
        )
      },
      {
        path: '/monitor',
        element: (
          <Suspense fallback={<CircularProgress />}>
            <LiveStream />
          </Suspense >
        )
      },
    ]
  },
])