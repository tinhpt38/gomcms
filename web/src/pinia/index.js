import { createPinia } from 'pinia'
import { useAppStore } from '@/pinia/modules/app'
import { useUserStore } from '@/pinia/modules/user'
import { useFileProcessStore } from '@/pinia/bl/fileProcess'
import { useDictionaryStore } from '@/pinia/modules/dictionary'

const store = createPinia()

export {
  store,
  useAppStore,
  useUserStore,
  useFileProcessStore,

  useDictionaryStore
}
