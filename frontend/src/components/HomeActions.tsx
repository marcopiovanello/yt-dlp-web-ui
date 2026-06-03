import { useSetAtom } from 'jotai'
import { Suspense, useState } from 'react'
import { loadingAtom } from '../atoms/ui'
import { useToast } from '../hooks/toast'
import DownloadDialog from './DownloadDialog'
import HomeSpeedDial from './HomeSpeedDial'
import TemplatesEditor from './TemplatesEditor'

const HomeActions: React.FC = () => {
  const setIsLoading = useSetAtom(loadingAtom)

  const [openDownload, setOpenDownload] = useState(false)
  const [openEditor, setOpenEditor] = useState(false)

  const { pushMessage } = useToast()

  return (
    <>
      <HomeSpeedDial
        onDownloadOpen={() => setOpenDownload(true)}
        onEditorOpen={() => setOpenEditor(true)}
      />
      <Suspense>
        <DownloadDialog
          open={openDownload}
          onClose={() => {
            setOpenDownload(false)
            setIsLoading(true)
          }}
          // TODO: handle optimistic UI update
          onDownloadStart={(url) => {
            pushMessage(`Requested ${url}`, 'info')
            setOpenDownload(false)
            setIsLoading(true)
          }}
        />
      </Suspense>
      <TemplatesEditor
        open={openEditor}
        onClose={() => setOpenEditor(false)}
      />
    </>
  )
}

export default HomeActions