export type MarkdownTable = {
  headers: string[]
  rows: string[][]
}

const MARKDOWN_TABLE_HEADER_RE = /(^|\n)\|\s*序号\s*\|/
const MARKDOWN_TABLE_DATA_ROW_RE = /(^|\n)\|\s*\d+\s*\|/

export function hasRenderableMarkdownTable(markdown: string): boolean {
  const normalized = normalizeMarkdown(markdown)
  return MARKDOWN_TABLE_HEADER_RE.test(normalized) && MARKDOWN_TABLE_DATA_ROW_RE.test(normalized)
}

export function parseMarkdownTable(markdown: string): MarkdownTable | null {
  if (!hasRenderableMarkdownTable(markdown)) return null

  const normalized = normalizeMarkdown(markdown)
  const lines = normalized
    .split('\n')
    .map((l) => l.trim())
    .filter((l) => l.includes('|'))

  if (lines.length < 3) return null

  const headerLine = lines[0]
  const sepLine = lines[1]
  if (!sepLine.includes('---')) return null

  let headers = splitRow(headerLine)
  let rows = lines.slice(2).map(splitRow).filter((r) => r.length > 0)
  if (headers.length === 0 || rows.length === 0) return null

  // 查找序号列的位置
  const indexCol = findColumnIndex(headers, ['序号'])
  
  // 如果序号列存在且不是第一列，删除序号列前面的所有列
  if (indexCol > 0) {
    headers = headers.slice(indexCol)
    rows = rows.map(row => row.slice(indexCol))
  }

  return { headers, rows }
}

function normalizeMarkdown(markdown: string): string {
  return markdown.replace(/\r\n/g, '\n').replace(/\n/g, '\n').replace(/\r\n/g, '\n')
}

function splitRow(line: string): string[] {
  const raw = line.split('|').map((c) => c.trim())
  const trimmed = raw.filter((c, idx) => !(idx === 0 && c === '') && !(idx === raw.length - 1 && c === ''))
  return trimmed
}

export type SignalRow = {
  index: number
  symbol: string
  direction: string
  entryRange: string
  stopLoss: string
  takeProfit: string
  holdingTime: string
  raw: Record<string, string>
}

function extractSignalsFromMarkdownLegacy(markdown: string): SignalRow[] {
  const table = parseMarkdownTable(markdown)
  if (!table) return []

  const headerIndex: Record<string, number> = {}
  table.headers.forEach((h, i) => {
    headerIndex[h] = i
  })

  const symbolCol = findColumnIndex(table.headers, ['品种', '品种（代码）', '品种(代码)'])
  const directionCol = findColumnIndex(table.headers, ['多空'])
  const entryCol = findColumnIndex(table.headers, ['入场区间'])
  const stopLossCol = findColumnIndex(table.headers, ['止损', '止损位'])
  const takeProfitCol = findColumnIndex(table.headers, ['止盈', '止盈位'])
  const holdingTimeCol = findColumnIndex(table.headers, ['持仓时间', '持有时间'])
  const indexCol = findColumnIndex(table.headers, ['序号'])

  if (symbolCol < 0 || directionCol < 0 || entryCol < 0) return []

  return table.rows
    .map((cols) => {
      const raw: Record<string, string> = {}
      table.headers.forEach((h, i) => {
        raw[h] = cols[i] ?? ''
      })

      const idxRaw = indexCol >= 0 ? cols[indexCol] : ''
      const idx = Number.parseInt(idxRaw, 10)

      return {
        index: Number.isFinite(idx) ? idx : 0,
        symbol: cols[symbolCol] ?? '',
        direction: cols[directionCol] ?? '',
        entryRange: cols[entryCol] ?? '',
        stopLoss: stopLossCol >= 0 ? cols[stopLossCol] ?? '' : '',
        takeProfit: takeProfitCol >= 0 ? cols[takeProfitCol] ?? '' : '',
        holdingTime: holdingTimeCol >= 0 ? cols[holdingTimeCol] ?? '' : '',
        raw,
      }
    })
    .filter((r) => r.symbol !== '')
}

export function extractSignalsFromMarkdown(markdown: string): SignalRow[] {
  const table = parseMarkdownTable(markdown)
  if (!table) return []

  const indexCol = findColumnIndex(table.headers, ['\u5e8f\u53f7'])
  const symbolCol = findColumnIndex(table.headers, [
    '\u54c1\u79cd',
    '\u54c1\u79cd\uff08\u4ee3\u7801\uff09',
    '\u54c1\u79cd(\u4ee3\u7801)',
    '\u671f\u6743\u5408\u7ea6',
  ])
  const directionCol = findColumnIndex(table.headers, [
    '\u591a\u7a7a',
    '\u7b56\u7565',
    '\u671f\u6743\u7b56\u7565',
  ])
  const entryCol = findColumnIndex(table.headers, [
    '\u5165\u573a\u533a\u95f4',
    '\u5165\u573a\u6743\u5229\u91d1\u533a\u95f4',
  ])
  const stopLossCol = findColumnIndex(table.headers, [
    '\u6b62\u635f',
    '\u6b62\u635f\u4f4d',
    '\u6b62\u635f\u6761\u4ef6',
  ])
  const takeProfitCol = findColumnIndex(table.headers, [
    '\u6b62\u76c8',
    '\u6b62\u76c8\u4f4d',
    '\u6b62\u76c8\u6761\u4ef6',
  ])
  const holdingTimeCol = findColumnIndex(table.headers, [
    '\u6301\u4ed3\u65f6\u95f4',
    '\u6301\u6709\u65f6\u95f4',
    '\u5efa\u8bae\u6301\u4ed3\u65f6\u95f4',
  ])

  if (symbolCol < 0 || directionCol < 0 || entryCol < 0) {
    return extractSignalsFromMarkdownLegacy(markdown)
  }

  return table.rows
    .map((cols) => {
      const raw: Record<string, string> = {}
      table.headers.forEach((h, i) => {
        raw[h] = cols[i] ?? ''
      })

      const idxRaw = indexCol >= 0 ? cols[indexCol] : ''
      const idx = Number.parseInt(idxRaw, 10)

      return {
        index: Number.isFinite(idx) ? idx : 0,
        symbol: cols[symbolCol] ?? '',
        direction: cols[directionCol] ?? '',
        entryRange: cols[entryCol] ?? '',
        stopLoss: stopLossCol >= 0 ? cols[stopLossCol] ?? '' : '',
        takeProfit: takeProfitCol >= 0 ? cols[takeProfitCol] ?? '' : '',
        holdingTime: holdingTimeCol >= 0 ? cols[holdingTimeCol] ?? '' : '',
        raw,
      }
    })
    .filter((r) => r.symbol !== '')
}

function findColumnIndex(headers: string[], candidates: string[]): number {
  for (const c of candidates) {
    const idx = headers.findIndex((h) => h.replace(/\s+/g, '') === c.replace(/\s+/g, ''))
    if (idx >= 0) return idx
  }
  for (const c of candidates) {
    const idx = headers.findIndex((h) => h.includes(c))
    if (idx >= 0) return idx
  }
  return -1
}
