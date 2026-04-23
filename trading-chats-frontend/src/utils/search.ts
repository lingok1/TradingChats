export function normalizeSearchKeyword(value: string): string {
  return value.trim().toLowerCase()
}

function stringifySearchValue(value: unknown): string {
  if (value == null) return ''
  if (Array.isArray(value)) {
    return value.map((item) => stringifySearchValue(item)).join(' ')
  }
  return String(value)
}

function escapeHtml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;')
}

export function matchesKeyword(values: unknown[], keyword: string): boolean {
  const normalizedKeyword = normalizeSearchKeyword(keyword)
  if (!normalizedKeyword) {
    return true
  }

  return values.some((value) => stringifySearchValue(value).toLowerCase().includes(normalizedKeyword))
}

export function highlightKeyword(value: unknown, keyword: string): string {
  const text = stringifySearchValue(value)
  const normalizedKeyword = normalizeSearchKeyword(keyword)
  if (!text || !normalizedKeyword) {
    return escapeHtml(text)
  }

  const lowerText = text.toLowerCase()
  const matchLength = normalizedKeyword.length
  let start = 0
  let result = ''

  while (start < text.length) {
    const index = lowerText.indexOf(normalizedKeyword, start)
    if (index === -1) {
      result += escapeHtml(text.slice(start))
      break
    }

    result += escapeHtml(text.slice(start, index))
    result += `<mark class="tc-search-highlight">${escapeHtml(text.slice(index, index + matchLength))}</mark>`
    start = index + matchLength
  }

  return result
}
