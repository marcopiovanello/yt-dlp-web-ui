import { useAtom, useAtomValue } from 'jotai'
import { useEffect } from 'react'
import { loadingDownloadsState } from '../atoms/downloads'
import { listViewState } from '../atoms/settings'
import { loadingAtom } from '../atoms/ui'
import DownloadsGridView from './DownloadsGridView'
import DownloadsTableView from './DownloadsTableView'

const Downloads: React.FC = () => {
  const tableView = useAtomValue(listViewState)
  const loadingDownloads = useAtomValue(loadingDownloadsState)

  const [isLoading, setIsLoading] = useAtom(loadingAtom)

  useEffect(() => {
    if (loadingDownloads) {
      return setIsLoading(true)
    }
    setIsLoading(false)
  }, [loadingDownloads, isLoading])

  if (tableView) return <DownloadsTableView />

  return <DownloadsGridView />
}

export default Downloads