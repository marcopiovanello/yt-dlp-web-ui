import {
  Container
} from '@mui/material'
import { loadingAtom } from '../atoms/ui'
import Downloads from '../components/Downloads'
import HomeActions from '../components/HomeActions'
import LoadingBackdrop from '../components/LoadingBackdrop'
import Splash from '../components/Splash'
import { useAtomValue } from 'jotai'

export default function Home() {
  const isLoading = useAtomValue(loadingAtom)

  return (
    <Container maxWidth="xl" sx={{ mt: 2, mb: 8 }}>
      <LoadingBackdrop isLoading={isLoading} />
      <Splash />
      <Downloads />
      <HomeActions />
    </Container>
  )
}
