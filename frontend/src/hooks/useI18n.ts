import Translator from '../lib/i18n'

export const useI18n = () => {
  const instance = Translator.instance

  return {
    i18n: instance,
    t: instance.t
  }
}