import type { ApiResponse } from './types'

export function unwrap<T>(res: ApiResponse<T>): T {
  if (res.code !== 200) {
    throw new Error(res.msg || `API error: ${res.code}`)
  }
  return res.data
}

