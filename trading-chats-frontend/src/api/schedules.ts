import { http } from './http'
import { unwrap } from './unwrap'
import type { ApiResponse, ScheduleConfig, ScheduleLog } from './types'

export async function getSchedules(): Promise<ScheduleConfig[]> {
  const res = await http.get<ApiResponse<ScheduleConfig[]>>('/schedules')
  return unwrap(res.data)
}

export async function createSchedule(body: Omit<ScheduleConfig, 'id' | 'created_at' | 'updated_at'>): Promise<ScheduleConfig> {
  const res = await http.post<ApiResponse<ScheduleConfig>>('/schedules', body)
  return unwrap(res.data)
}

export async function deleteSchedule(id: string): Promise<string> {
  const res = await http.delete<ApiResponse<string>>(`/schedules/${id}`)
  return unwrap(res.data)
}

export async function updateScheduleStatus(id: string, status: 'active' | 'paused'): Promise<string> {
  const res = await http.put<ApiResponse<string>>('/schedules/status', { id, status })
  return unwrap(res.data)
}

export async function getScheduleLogs(id: string): Promise<ScheduleLog[]> {
  const res = await http.get<ApiResponse<ScheduleLog[]>>(`/schedules/${id}/logs`)
  return unwrap(res.data)
}

