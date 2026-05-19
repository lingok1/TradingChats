import { http } from './http'
import type { ApiResponse, TenantConfig, TenantMenuConfig } from './types'
import { unwrap } from './unwrap'

export type Tenant = {
  id: string
  name: string
  code: string
  type: string
  status: string
}

export async function getTenants(): Promise<Tenant[]> {
  const res = await http.get<ApiResponse<Tenant[]>>('/auth/tenants')
  return unwrap(res.data)
}

export async function getTenantConfig(tenantId?: string): Promise<TenantConfig> {
  const params = tenantId ? { tenant_id: tenantId } : {}
  const res = await http.get<ApiResponse<TenantConfig>>('/tenant-config', { params })
  return unwrap(res.data)
}

export async function saveTenantMenu(tenantId: string, menu: TenantMenuConfig): Promise<void> {
  const res = await http.put<ApiResponse<void>>(`/tenant-config/menu?tenant_id=${tenantId}`, menu)
  unwrap(res.data)
}

export async function saveTenantParameters(
  params: Record<string, string>,
  tenantId?: string,
): Promise<void> {
  const query = tenantId ? `?tenant_id=${tenantId}` : ''
  const res = await http.put<ApiResponse<void>>(`/tenant-config/parameters${query}`, { parameters: params })
  unwrap(res.data)
}

export async function listTenantConfigs(): Promise<TenantConfig[]> {
  const res = await http.get<ApiResponse<TenantConfig[]>>('/tenant-config/list')
  return unwrap(res.data)
}
