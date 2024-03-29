import i18n, { Resource } from 'i18next';
import { initReactI18next } from 'react-i18next';

import enUS from './en-US';
import zhCN from './zh-CN';

// Ready translated locale messages
export const langs: { [key: string]: string } = {
  'en-US': 'English',
  'zh-CN': '简体中文',
};

// Ready translated locale messages
export const resources: Resource = {
  'en-US': {
    translation: enUS,
  },
  'zh-CN': {
    translation: zhCN,
  },
};

export function changeLanguage(lang: string): void {
  localStorage.setItem('locale', lang);
  i18n.changeLanguage(lang);
}

i18n
  .use(initReactI18next) // passes i18n down to react-i18next
  .init({
    resources,
    lng: 'zh-CN',
    fallbackLng: 'en-US',

    interpolation: {
      escapeValue: false,
    },
  });

const locale = localStorage.getItem('locale');
if (locale) {
  i18n.changeLanguage(locale);
}
