import AddCircleIcon from '@mui/icons-material/AddCircle'
import BuildCircleIcon from '@mui/icons-material/BuildCircle'
import ClearAllIcon from '@mui/icons-material/ClearAll'
import DeleteForeverIcon from '@mui/icons-material/DeleteForever'
import FolderZipIcon from '@mui/icons-material/FolderZip'
import FormatListBulleted from '@mui/icons-material/FormatListBulleted'
import ViewAgendaIcon from '@mui/icons-material/ViewAgenda'
import {
  SpeedDial,
  SpeedDialAction,
  SpeedDialIcon
} from '@mui/material'
import { useAtom, useAtomValue } from 'jotai'
import { listViewState, serverURL } from '../atoms/settings'
import { useI18n } from '../hooks/useI18n'
import { useRPC } from '../hooks/useRPC'

type Props = {
  onDownloadOpen: () => void
  onEditorOpen: () => void
}

const HomeSpeedDial: React.FC<Props> = ({ onDownloadOpen, onEditorOpen }) => {
  const serverAddr = useAtomValue(serverURL)
  const [listView, setListView] = useAtom(listViewState)

  const { i18n } = useI18n()
  const { client } = useRPC()

  return (
    <SpeedDial
      ariaLabel="Home speed dial"
      sx={{ position: 'absolute', bottom: 64, right: 24 }}
      icon={<SpeedDialIcon />}
    >
      <SpeedDialAction
        icon={listView ? <ViewAgendaIcon /> : <FormatListBulleted />}
        tooltipTitle={listView ? 'Card view' : i18n.t('tableView')}
        onClick={() => setListView(state => !state)}
      />
      <SpeedDialAction
        icon={<FolderZipIcon />}
        tooltipTitle={i18n.t('bulkDownload')}
        onClick={() => window.open(`${serverAddr}/archive/bulk?token=${localStorage.getItem('token')}`)}
      />
      <SpeedDialAction
        icon={<ClearAllIcon />}
        tooltipTitle={i18n.t('clearCompletedButton')}
        onClick={() => client.clearCompleted()}
      />
      <SpeedDialAction
        icon={<DeleteForeverIcon />}
        tooltipTitle={i18n.t('abortAllButton')}
        onClick={() => client.killAll()}
      />
      <SpeedDialAction
        icon={<BuildCircleIcon />}
        tooltipTitle={i18n.t('templatesEditor')}
        onClick={onEditorOpen}
      />
      <SpeedDialAction
        icon={<AddCircleIcon />}
        tooltipTitle={i18n.t('newDownloadButton')}
        onClick={onDownloadOpen}
      />
    </SpeedDial>
  )
}

export default HomeSpeedDial