export type ApiResponse<T> = {
  code: number
  msg: string
  data: T
}

export type TimeLike = string | number | null | undefined | Record<string, unknown>

export type AIResponse = {
  id?: string
  batch_id: string
  prompt?: string
  response: string
  model_name: string
  provider: string
  status: string
  error?: string
  created_at?: TimeLike
  updated_at?: TimeLike
  completed_at?: TimeLike
}

export type GenerateAIRequest = {
  template_id: string
  param1: string
  param2: string
}

export type GenerateAIResponse = {
  batch_id: string
}

export type PromptTemplate = {
  id?: string
  name: string
  content: string
  tags: string[]
  created_at?: TimeLike
  updated_at?: TimeLike
}

export type ModelAPIConfig = {
  id?: string
  name: string
  api_url: string
  api_key: string
  models: string[]
  provider: string
  created_at?: TimeLike
  updated_at?: TimeLike
}

export type ScheduleConfig = {
  id?: string
  name: string
  cron_expr: string
  template_id: string
  status: 'active' | 'paused'
  created_at?: TimeLike
  updated_at?: TimeLike
}

export type ScheduleLog = {
  id?: string
  schedule_config_id: string
  batch_id: string
  status: string
  error?: string
  executed_at?: TimeLike
}

export type SystemConfig = {
  id?: string
  system_title: string
  system_logo: string
  parameters: Record<string, string>
  updated_at?: TimeLike
}

export type NewsItem = {
  id: string
  title: string
  summary: string
  content: string
  category: string
  source: string
  author?: string
  publish_time: string
  image_url?: string
  read_count?: number
}

export type NewsCategory = {
  id: string
  name: string
  code: string
  icon?: string
}
