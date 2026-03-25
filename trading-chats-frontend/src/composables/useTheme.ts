import { computed, onMounted, ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'

const storageKey = 'trading-chats-theme'

export function useTheme() {
  const mode = ref<ThemeMode>('system')
  const systemDark = ref(false)

  const isDark = computed(() => {
    if (mode.value === 'dark') return true
    if (mode.value === 'light') return false
    return systemDark.value
  })

  function apply() {
    const root = document.documentElement
    if (isDark.value) root.classList.add('dark')
    else root.classList.remove('dark')
  }

  function load() {
    const v = localStorage.getItem(storageKey)
    if (v === 'light' || v === 'dark' || v === 'system') mode.value = v
  }

  function save() {
    localStorage.setItem(storageKey, mode.value)
  }

  onMounted(() => {
    load()

    const mq = window.matchMedia('(prefers-color-scheme: dark)')
    systemDark.value = mq.matches
    mq.addEventListener('change', (e) => {
      systemDark.value = e.matches
    })

    apply()
  })

  watch([mode, systemDark], () => {
    save()
    apply()
  })

  return { mode, isDark }
}

