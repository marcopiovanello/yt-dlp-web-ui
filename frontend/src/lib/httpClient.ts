import { tryCatch } from 'fp-ts/TaskEither'
import * as J from 'fp-ts/Json'
import * as E from 'fp-ts/Either'
import { pipe } from 'fp-ts/lib/function'

async function fetcher(url: string, opt?: RequestInit, controller?: AbortController): Promise<string> {
  const jwt = localStorage.getItem('token')

  if (opt && !opt.headers) {
    opt.headers = {
      'Content-Type': 'application/json',
    }
  }

  const res = await fetch(url, {
    ...opt,
    headers: {
      ...opt?.headers,
      'X-Authentication': jwt ?? ''
    },
    signal: controller?.signal
  })

  if (!res.ok) {
    throw await res.text()
  }



  return res.text()
}

export const ffetch = <T>(url: string, opt?: RequestInit, controller?: AbortController) => tryCatch(
  async () => pipe(
    await fetcher(url, opt, controller),
    J.parse,
    E.match(
      (l) => l as T,
      (r) => r as T
    )
  ),
  (e) => `error while fetching: ${e}`
)
