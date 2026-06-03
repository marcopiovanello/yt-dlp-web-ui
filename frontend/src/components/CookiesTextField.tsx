import { Button, TextField } from '@mui/material'
import * as A from 'fp-ts/Array'
import * as E from 'fp-ts/Either'
import * as O from 'fp-ts/Option'
import { matchW } from 'fp-ts/lib/TaskEither'
import { pipe } from 'fp-ts/lib/function'
import { useMemo } from 'react'
import { Subject, debounceTime, distinctUntilChanged } from 'rxjs'
import { serverSideCookiesState, serverURL } from '../atoms/settings'
import { useSubscription } from '../hooks/observable'
import { useToast } from '../hooks/toast'
import { ffetch } from '../lib/httpClient'
import { useAtomValue } from 'jotai'
import { useI18n } from '../hooks/useI18n'

const { i18n } = useI18n()

const validateCookie = (cookie: string) => pipe(
  cookie,
  cookie => cookie.split('\t'),
  E.of,
  E.flatMap(
    E.fromPredicate(
      f => f.length === 7,
      () => `missing parts`
    )
  ),
  E.flatMap(
    E.fromPredicate(
      f => f[0].length > 0,
      () => 'missing domain'
    )
  ),
  E.flatMap(
    E.fromPredicate(
      f => f[1] === 'TRUE' || f[1] === 'FALSE',
      () => `invalid include subdomains`
    )
  ),
  E.flatMap(
    E.fromPredicate(
      f => f[2].length > 0,
      () => 'invalid path'
    )
  ),
  E.flatMap(
    E.fromPredicate(
      f => f[3] === 'TRUE' || f[3] === 'FALSE',
      () => 'invalid secure flag'
    )
  ),
  E.flatMap(
    E.fromPredicate(
      f => isFinite(Number(f[4])),
      () => 'invalid expiration'
    )
  ),
  E.flatMap(
    E.fromPredicate(
      f => f[5].length > 0,
      () => 'invalid name'
    )
  ),
  E.flatMap(
    E.fromPredicate(
      f => f[6].length > 0,
      () => 'invalid value'
    )
  ),
)

const noopValidator = (s: string): E.Either<string, string[]> => pipe(
  s,
  s => s.split('\t'),
  E.of
)

const isCommentOrNewLine = (s: string) => s === '' || s.startsWith('\n') || s.startsWith('#')

const CookiesTextField: React.FC = () => {
  const serverAddr = useAtomValue(serverURL)
  const savedCookies = useAtomValue(serverSideCookiesState)

  const { pushMessage } = useToast()

  const cookies$ = useMemo(() => new Subject<string>(), [])

  const submitCookies = (cookies: string) =>
    ffetch(`${serverAddr}/api/v1/cookies`, {
      method: 'POST',
      body: JSON.stringify({
        cookies
      })
    })()

  const deleteCookies = () => pipe(
    ffetch(`${serverAddr}/api/v1/cookies`, {
      method: 'DELETE',
    }),
    matchW(
      (l) => pushMessage(l, 'error'),
      (_) => {
        pushMessage('Deleted cookies', 'success')
        pushMessage(`Reload the page to apply the changes`, 'info')
      }
    )
  )()

  const validateNetscapeCookies = (cookies: string) => pipe(
    cookies,
    cookies => cookies.split('\n'),
    A.map(c => isCommentOrNewLine(c) ? noopValidator(c) : validateCookie(c)), // validate line
    A.mapWithIndex((i, either) => pipe(                  // detect errors and return the either
      either,
      E.match(
        (l) => {
          pushMessage(`Error in line ${i + 1}: ${l}`, 'warning')
          return either
        },
        (_) => either
      ),
    )),
    A.filter(c => E.isRight(c)),                         // filter the line who didn't pass the validation
    A.map(E.getOrElse(() => new Array<string>())),       // cast the array of eithers to an array of tokens
    A.filter(f => f.length > 0),                         // filter the empty tokens
    A.map(f => f.join('\t')),                            // join the tokens in a TAB separated string
    A.reduce('', (c, n) => `${c}${n}\n`),                // reduce all to a single string separated by \n
    parsed => parsed.length > 0                          // if nothing has passed the validation return none
      ? O.some(parsed)
      : O.none
  )

  useSubscription(
    cookies$.pipe(
      debounceTime(650),
      distinctUntilChanged()
    ),
    (cookies) => pipe(
      cookies,
      validateNetscapeCookies,
      O.match(
        () => pushMessage('No valid cookies', 'warning'),
        async (some) => {
          pipe(
            await submitCookies(some.trimEnd()),
            E.match(
              (l) => pushMessage(`${l}`, 'error'),
              () => {
                pushMessage(`Saved ${some.split('\n').length} Netscape cookies`, 'success')
                pushMessage('Reload the page to apply the changes', 'info')
              }
            )
          )
        }
      )
    )
  )

  return (
    <>
      <TextField
        label="Netscape Cookies"
        multiline
        maxRows={20}
        minRows={4}
        fullWidth
        defaultValue={savedCookies}
        onChange={(e) => cookies$.next(e.currentTarget.value)}
      />
      <Button onClick={deleteCookies}>{i18n.t('deleteCookies')}</Button>
    </>
  )
}

export default CookiesTextField
