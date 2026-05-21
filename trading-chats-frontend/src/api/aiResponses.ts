import { http } from './http'
import type { AIResponse, ApiResponse, GenerateAIRequest, GenerateAIResponse, TabTag } from './types'
import { unwrap } from './unwrap'

export async function getLatestAIResponses(tabTag: TabTag = 'futures', subTag?: string): Promise<AIResponse[]> {
  const params: Record<string, string> = { tab_tag: tabTag }
  if (subTag) params.sub_tag = subTag
  const res = await http.get<ApiResponse<AIResponse[]>>('/ai-responses/latest', { params })
  return unwrap(res.data).filter((item) => item.status === 'completed')
}

export async function getAIResponsesByBatch(batchId: string, tabTag: TabTag = 'futures'): Promise<AIResponse[]> {
  const res = await http.get<ApiResponse<AIResponse[]>>('/ai-responses/batch', {
    params: { batch_id: batchId, tab_tag: tabTag },
  })
  return unwrap(res.data)
}

export async function generateAIResponses(req: GenerateAIRequest): Promise<GenerateAIResponse> {
  const res = await http.post<ApiResponse<GenerateAIResponse>>('/ai-responses/generate', req)
  return unwrap(res.data)
}

export function getAIResponseEventsUrl(tabTag: string = 'futures'): string {
  const url = new URL('/api/ai-responses/events', window.location.origin)
  url.searchParams.set('tab_tag', tabTag)
  return url.toString()
}
