import { http } from './http'
import { unwrap } from './unwrap'
import type { ApiResponse } from './types'

export type LoginRequest = {
  username: string
  password: string
}

export type LoginResponse = {
  access_token: string
  refresh_token: string
  token_type: string
  expires_at: string
  user: {
    id: string
    tenant_id: string
    username: string
    display_name: string
    role: string
    status: string
    last_login_at?: string
    created_at?: string
    updated_at?: string
  }
}

export async function login(body: LoginRequest): Promise<LoginResponse> {
  const res = await http.post<ApiResponse<LoginResponse>>('/auth/login', body)
  return unwrap(res.data)
}
