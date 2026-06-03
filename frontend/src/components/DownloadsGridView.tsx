import { Grid2 } from '@mui/material'
import { useAtomValue } from 'jotai'
import { useTransition } from 'react'
import { activeDownloadsState } from '../atoms/downloads'
import { useToast } from '../hooks/toast'
import { useI18n } from '../hooks/useI18n'
import { useRPC } from '../hooks/useRPC'
import { ProcessStatus, RPCResult } from '../types'
import DownloadCard from './DownloadCard'
import LoadingBackdrop from './LoadingBackdrop'

const DownloadsGridView: React.FC = () => {
  const downloads = useAtomValue(activeDownloadsState)

  const { i18n } = useI18n()
  const { client } = useRPC()
  const { pushMessage } = useToast()

  const [isPending, startTransition] = useTransition()

  const stop = async (r: RPCResult) => r.progress.process_status === ProcessStatus.COMPLETED
    ? await client.clear(r.id)
    : await client.kill(r.id)

  return (
    <>
      <LoadingBackdrop isLoading={isPending} />
      <Grid2 container spacing={{ xs: 2, md: 2 }} columns={{ xs: 4, sm: 8, md: 12, xl: 12 }} pt={2}>
        {
          downloads.map(download => (
            <Grid2 size={{ xs: 4, sm: 8, md: 6, xl: 4 }} key={download.id}>
              <DownloadCard
                download={download}
                onStop={() => startTransition(async () => {
                  await stop(download)
                })}
                onCopy={() => pushMessage(i18n.t('clipboardAction'), 'info')}
              />
            </Grid2>
          ))
        }
      </Grid2>
    </>
  )
}

export default DownloadsGridView