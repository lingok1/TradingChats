import { http } from './http'
import type { AIResponse, ApiResponse, GenerateAIRequest, GenerateAIResponse } from './types'
import { unwrap } from './unwrap'

export async function getLatestAIResponses(): Promise<AIResponse[]> {
  const res = await http.get<ApiResponse<AIResponse[]>>('/ai-responses/latest')
  return unwrap(res.data)
}

export async function getAIResponsesByBatch(batchId: string): Promise<AIResponse[]> {
  const res = await http.get<ApiResponse<AIResponse[]>>('/ai-responses/batch', {
    params: { batch_id: batchId },
  })
  return unwrap(res.data)
}

export async function generateAIResponses(req: GenerateAIRequest): Promise<GenerateAIResponse> {
  const res = await http.post<ApiResponse<GenerateAIResponse>>('/ai-responses/generate', req)
  return unwrap(res.data)
}
