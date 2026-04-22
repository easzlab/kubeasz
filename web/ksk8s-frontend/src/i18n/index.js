import { createI18n } from 'vue-i18n'
import en from '../locales/en.json'
import zhCN from '../locales/zh-CN.json'
import zhTW from '../locales/zh-TW.json'
import fr from '../locales/fr.json'

const savedLocale = localStorage.getItem('ksk8s_locale')
const browserLocale = navigator.language

function detectLocale() {
  if (savedLocale) return savedLocale
  if (browserLocale.startsWith('zh')) {
    if (browserLocale === 'zh-TW' || browserLocale === 'zh-HK') return 'zh-TW'
    return 'zh-CN'
  }
  if (browserLocale.startsWith('fr')) return 'fr'
  return 'en'
}

const i18n = createI18n({
  legacy: false,
  locale: detectLocale(),
  fallbackLocale: 'en',
  messages: {
    en,
    'zh-CN': zhCN,
    'zh-TW': zhTW,
    fr
  }
})

export default i18n
