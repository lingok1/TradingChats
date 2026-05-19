import { http } from './http'
import type { ApiResponse } from './types'
import { unwrap } from './unwrap'

export type RecommendationItem = {
  symbol: string
  direction: string
  entry_range: string
  take_profit: string
  stop_loss: string
  reason: string
}

export type FuturesRecommendation = {
  id: string
  batch_id: string
  items: RecommendationItem[]
  model_name: string
  model_api_name: string
  created_at: string
}

export async function getLatestFuturesRecommendation(): Promise<FuturesRecommendation | null> {
  const res = await http.get<ApiResponse<FuturesRecommendation | null>>('/futures-recommendation/latest')
  return unwrap(res.data)
}
