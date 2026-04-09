import { http } from './http'
import { unwrap } from './unwrap'
import type { ApiResponse, PromptTemplate } from './types'

export async function getPromptTemplates(): Promise<PromptTemplate[]> {
  const res = await http.get<ApiResponse<PromptTemplate[]>>('/prompt-templates')
  return unwrap(res.data)
}

export async function createPromptTemplate(body: Omit<PromptTemplate, 'id' | 'created_at' | 'updated_at'>): Promise<PromptTemplate> {
  const res = await http.post<ApiResponse<PromptTemplate>>('/prompt-templates', body)
  return unwrap(res.data)
}

export async function updatePromptTemplate(
  id: string,
  body: Omit<PromptTemplate, 'id' | 'created_at' | 'updated_at'>,
): Promise<PromptTemplate> {
  const res = await http.put<ApiResponse<PromptTemplate>>(`/prompt-templates/${id}`, body)
  return unwrap(res.data)
}

export async function deletePromptTemplate(id: string): Promise<{message: string}> {
  const res = await http.delete<ApiResponse<{message: string}>>(`/prompt-templates/${id}`)
  return unwrap(res.data)
}

export async function generatePrompt(body: {
  template_id: string
  param1?: string
  param2?: string
}): Promise<string> {
  const res = await http.post<ApiResponse<{ prompt: string }>>('/prompt-templates/generate', body)
  return unwrap(res.data).prompt
}

