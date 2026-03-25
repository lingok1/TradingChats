import { http } from './http'
import type { ApiResponse, SystemConfig } from './types'
import { unwrap } from './unwrap'

export type SaveSystemBasicConfigRequest = Pick<SystemConfig, 'system_title' | 'system_logo'>
export type SaveSystemParametersRequest = Pick<SystemConfig, 'parameters'>

export async function getSystemConfig(): Promise<SystemConfig> {
  const res = await http.get<ApiResponse<SystemConfig>>('/system-config')
  return unwrap(res.data)
}

export async function saveSystemBasicConfig(config: SaveSystemBasicConfigRequest): Promise<void> {
  const res = await http.put<ApiResponse<void>>('/system-config/basic', config)
  unwrap(res.data)
}

export async function saveSystemParameters(config: SaveSystemParametersRequest): Promise<void> {
  const res = await http.put<ApiResponse<void>>('/system-config/parameters', config)
  unwrap(res.data)
}
