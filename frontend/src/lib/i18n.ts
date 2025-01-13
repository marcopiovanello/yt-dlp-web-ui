//@ts-ignore
import i18n from '../assets/i18n.yaml'
//@ts-ignore
import fallback from '../assets/i18n/en_US.yaml'

export default class Translator {
  static #instance: Translator

  private language: string
  private current: string[] = []

  constructor() {
    this.language = localStorage.getItem('language')?.replaceAll('"', '') ?? 'english'
    this.setLanguage(this.language)
  }

  getLanguage(): string {
    return this.language
  }

  async setLanguage(language: string): Promise<void> {
    this.language = language

    let isoCodeFile: string = i18n.languages[language]

    // extension needs to be in source code to help vite bundle all yaml files
    if (isoCodeFile.endsWith('.yaml')) {
      isoCodeFile = isoCodeFile.replaceAll('.yaml', '')
    }

    if (isoCodeFile) {
      const { default: translations } = await import(`../assets/i18n/${isoCodeFile}.yaml`)

      this.current = translations.keys
    }
  }

  t(key: string): string {
    if (this.current) {
      //@ts-ignore
      return this.current[key] ?? fallback.keys[key]
    }
    return 'caption not defined'
  }

  public static get instance(): Translator {
    if (!Translator.#instance) {
      Translator.#instance = new Translator()
    }

    return Translator.#instance
  }
}
