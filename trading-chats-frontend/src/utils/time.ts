import type { TimeLike } from '../api/types'

export function asTimeString(v: TimeLike): string {
  if (typeof v === 'string') return v
  if (typeof v === 'number') return String(v)
  if (!v) return ''
  if (typeof v === 'object') {
    const s = JSON.stringify(v)
    return s === '{}' ? '' : s
  }
  return ''
}
