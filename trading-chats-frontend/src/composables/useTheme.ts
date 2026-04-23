import { computed, onMounted, ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark'

const storageKey = 'trading-chats-theme'

export function useTheme() {
  const mode = ref<ThemeMode>('light')
  const isDark = computed(() => mode.value === 'dark')

  function apply() {
    const root = document.documentElement
    if (isDark.value) root.classList.add('dark')
    else root.classList.remove('dark')
  }

  function load() {
    const v = localStorage.getItem(storageKey)
    mode.value = v === 'dark' ? 'dark' : 'light'
  }

  function save() {
    localStorage.setItem(storageKey, mode.value)
  }

  onMounted(() => {
    load()

    apply()
  })

  watch(mode, () => {
    save()
    apply()
  })

  return { mode, isDark }
}

