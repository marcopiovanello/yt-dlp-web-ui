import { Alert, Snackbar } from "@mui/material"
import { Toast, toastListState } from '../atoms/toast'
import { useEffect } from 'react'
import { useAtom } from 'jotai'

const Toaster: React.FC = () => {
  const [toasts, setToasts] = useAtom(toastListState)

  const deletePredicate = (t: Toast) => (Date.now() - t.createdAt) < 2000

  useEffect(() => {
    if (toasts.length > 0) {
      const closer = setInterval(() => {
        setToasts(t => t.map(t => ({ ...t, open: deletePredicate(t) })))
      }, 900)

      const cleaner = setInterval(() => {
        setToasts(t => t.filter(deletePredicate))
      }, 2005)

      return () => {
        clearInterval(closer)
        clearInterval(cleaner)
      }
    }
  }, [setToasts, toasts.length])

  return (
    <>
      {toasts.map((toast, index) => (
        <Snackbar
          key={index}
          open={toast.open}
          sx={index > 0 ? { marginBottom: index * 6.5 } : null}
        >
          <Alert variant="filled" severity={toast.severity}>
            {toast.message}
          </Alert>
        </Snackbar>
      ))}
    </>
  )
}

export default Toaster