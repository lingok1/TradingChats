import { http } from './http'
import { unwrap } from './unwrap'
import type { ApiResponse, ModelAPIConfig } from './types'

export async function getModelApiConfigs(): Promise<ModelAPIConfig[]> {
  const res = await http.get<ApiResponse<ModelAPIConfig[]>>('/model-api-configs')
  return unwrap(res.data)
}

export async function createModelApiConfig(
  body: Omit<ModelAPIConfig, 'id' | 'created_at' | 'updated_at'>,
): Promise<ModelAPIConfig> {
  const res = await http.post<ApiResponse<ModelAPIConfig>>('/model-api-configs', body)
  return unwrap(res.data)
}

export async function updateModelApiConfig(
  id: string,
  body: Omit<ModelAPIConfig, 'id' | 'created_at' | 'updated_at'>,
): Promise<ModelAPIConfig> {
  const res = await http.put<ApiResponse<ModelAPIConfig>>(`/model-api-configs/${id}`, body)
  return unwrap(res.data)
}

export async function deleteModelApiConfig(id: string): Promise<string> {
  const res = await http.delete<ApiResponse<string>>(`/model-api-configs/${id}`)
  return unwrap(res.data)
}

export async function testModelApiConfig(id: string): Promise<{message: string}> {
  const res = await http.post<ApiResponse<{message: string}>>(`/model-api-configs/${id}/test`)
  return unwrap(res.data)
}

