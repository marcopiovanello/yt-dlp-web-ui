import { Container, FormControl, Grid2, InputLabel, MenuItem, Pagination, Select } from '@mui/material'
import { pipe } from 'fp-ts/lib/function'
import { matchW } from 'fp-ts/lib/TaskEither'
import { useAtomValue } from 'jotai'
import { useEffect, useState, useTransition } from 'react'
import { serverURL } from '../atoms/settings'
import ArchiveCard from '../components/ArchiveCard'
import EmptyArchive from '../components/EmptyArchive'
import { useToast } from '../hooks/toast'
import { ffetch } from '../lib/httpClient'
import { ArchiveEntry, PaginatedResponse } from '../types'
import LoadingBackdrop from '../components/LoadingBackdrop'

const Archive: React.FC = () => {
  const [isLoading, setLoading] = useState(true)
  const [archiveEntries, setArchiveEntries] = useState<ArchiveEntry[]>()

  const [currentCursor, setCurrentCursor] = useState(0)
  const [cursor, setCursor] = useState({ first: 0, next: 0 })
  const [pageSize, setPageSize] = useState(25)

  const [isPending, startTransition] = useTransition()

  const serverAddr = useAtomValue(serverURL)
  const { pushMessage } = useToast()

  const fetchArchived = (startCursor = 0) => pipe(
    ffetch<PaginatedResponse<ArchiveEntry[]>>(`${serverAddr}/archive?id=${startCursor}&limit=${pageSize}`),
    matchW(
      (l) => pushMessage(l, 'error'),
      (r) => {
        setArchiveEntries(r.data)
        setCursor({ ...cursor, first: r.first, next: r.next })
      }
    )
  )()

  const softDelete = (id: string) => pipe(
    ffetch<ArchiveEntry[]>(`${serverAddr}/archive/soft/${id}`, {
      method: 'DELETE'
    }),
    matchW(
      (l) => pushMessage(l, 'error'),
      (_) => startTransition(async () => await fetchArchived())
    )
  )()

  const hardDelete = (id: string) => pipe(
    ffetch<ArchiveEntry[]>(`${serverAddr}/archive/hard/${id}`, {
      method: 'DELETE'
    }),
    matchW(
      (l) => pushMessage(l, 'error'),
      (_) => startTransition(async () => await fetchArchived())
    )
  )()

  const setPage = (page: number) => setCurrentCursor(pageSize * (page - 1))

  useEffect(() => {
    fetchArchived(currentCursor).then(() => setLoading(false))
  }, [currentCursor])

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 8, minHeight: '80vh' }}>
      <LoadingBackdrop isLoading={isPending || isLoading} />
      {
        archiveEntries && archiveEntries.length !== 0 ?
          <Grid2
            container
            spacing={{ xs: 2, md: 2 }}
            columns={{ xs: 4, sm: 8, md: 12, xl: 12 }}
            pt={2}
            sx={{ minHeight: '77.5vh' }}
          >
            {
              archiveEntries.map((entry) => (
                <Grid2 size={{ xs: 4, sm: 8, md: 4, xl: 3 }} key={entry.id}>
                  <ArchiveCard
                    entry={entry}
                    onDelete={() => startTransition(async () => await softDelete(entry.id))}
                    onHardDelete={() => startTransition(async () => await hardDelete(entry.id))}
                  />
                </Grid2>
              ))
            }
          </Grid2>
          : <EmptyArchive />
      }
      <Pagination
        sx={{ mx: 'auto', pt: 2 }}
        count={Math.floor(cursor.next / pageSize) + 1}
        onChange={(_, v) => setPage(v)}
      />
      <FormControl variant="standard" sx={{ m: 1, minWidth: 120 }}>
        <InputLabel id="page-size-select-label">Page size</InputLabel>
        <Select
          labelId="page-size-select-label"
          value={pageSize}
          onChange={(e) => setPageSize(Number(e.target.value))}
          label="Page size"
        >
          <MenuItem value={25}>25</MenuItem>
          <MenuItem value={50}>50</MenuItem>
          <MenuItem value={100}>100</MenuItem>
        </Select>
      </FormControl>
    </Container>
  )
}

export default Archive