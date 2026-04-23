import { http } from './http'
import { unwrap } from './unwrap'
import type { ApiResponse, TradePlan } from './types'

type TradePlanPayload = Omit<TradePlan, 'id' | 'tenant_id' | 'created_at' | 'updated_at'>

export async function getTradePlans(tabTag: TradePlan['tab_tag']): Promise<TradePlan[]> {
  const res = await http.get<ApiResponse<TradePlan[]>>('/trade-plans', {
    params: { tab_tag: tabTag },
  })
  return unwrap(res.data)
}

export async function getTradePlanById(id: string): Promise<TradePlan> {
  const res = await http.get<ApiResponse<TradePlan>>(`/trade-plans/${id}`)
  return unwrap(res.data)
}

export async function createTradePlan(body: TradePlanPayload): Promise<TradePlan> {
  const res = await http.post<ApiResponse<TradePlan>>('/trade-plans', body)
  return unwrap(res.data)
}

export async function updateTradePlan(id: string, body: TradePlanPayload): Promise<string> {
  const res = await http.put<ApiResponse<string>>(`/trade-plans/${id}`, body)
  return unwrap(res.data)
}

export async function deleteTradePlan(id: string): Promise<string> {
  const res = await http.delete<ApiResponse<string>>(`/trade-plans/${id}`)
  return unwrap(res.data)
}
