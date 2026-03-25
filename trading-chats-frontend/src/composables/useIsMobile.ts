import { computed, onMounted, onUnmounted, ref } from 'vue'

export function useIsMobile(breakpointPx = 768) {
  const width = ref<number>(1024)
  const forced = ref<boolean>(false)

  function onResize() {
    width.value = window.innerWidth
  }

  onMounted(() => {
    const q = new URLSearchParams(window.location.search).get('mobile')
    forced.value = q === '1' || q === 'true'
    onResize()
    window.addEventListener('resize', onResize)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', onResize)
  })

  const isMobile = computed(() => forced.value || width.value < breakpointPx)
  return { isMobile, width }
}

