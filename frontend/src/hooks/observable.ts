import { useEffect, useState } from 'react'
import { Observable } from 'rxjs'

/**
 * Handles the subscription and unsubscription from an observable.
 * Automatically disposes the subscription.
 * @param source$ source observable
 * @param nextHandler subscriber function
 * @param errHandler error catching callback
 */
export function useSubscription<T>(
  source$: Observable<T>,
  nextHandler: (value: T) => void,
  errHandler?: (err: any) => void,
) {
  useEffect(() => {
    if (source$) {
      const sub = source$.subscribe({
        next: nextHandler,
        error: errHandler,
      })
      return () => sub.unsubscribe()
    }
  }, [source$])
}

/**
 * Use an observable as state
 * @param source$ source observable
 * @param initialState the initial state prior to the emission
 * @param errHandler error catching callback
 * @returns value emitted to the observable
 */
export function useObservable<T>(
  source$: Observable<T>,
  initialState: T,
  errHandler?: (err: any) => void,
): T {
  const [value, setValue] = useState(initialState)

  useSubscription(source$, setValue, errHandler)

  return value
}