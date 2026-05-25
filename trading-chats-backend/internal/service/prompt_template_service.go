package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromptTemplateService struct {
	repo                *repository.PromptTemplateRepository
	systemConfigService SystemConfigService
	tenantConfigService *TenantConfigService
	aiResponseRepo      *repository.AIResponseRepository
}

func NewPromptTemplateService(repo *repository.PromptTemplateRepository, systemConfigService SystemConfigService, tenantConfigService *TenantConfigService) *PromptTemplateService {
	return &PromptTemplateService{
		repo:                repo,
		systemConfigService: systemConfigService,
		tenantConfigService: tenantConfigService,
	}
}

// SetAIResponseRepo 注入 AIResponseRepo（用于推荐任务的 {{.xxxAIResult}} 占位符）
func (s *PromptTemplateService) SetAIResponseRepo(repo *repository.AIResponseRepository) {
	s.aiResponseRepo = repo
}

func (s *PromptTemplateService) CreatePromptTemplate(ctx context.Context, template *models.PromptTemplate) error {
	return s.repo.Create(ctx, template)
}

func (s *PromptTemplateService) GetPromptTemplateByID(ctx context.Context, id string) (*models.PromptTemplate, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.GetByID(ctx, objectID)
}

func (s *PromptTemplateService) GetAllPromptTemplates(ctx context.Context) ([]models.PromptTemplate, error) {
	return s.repo.GetAll(ctx)
}

func (s *PromptTemplateService) GetPromptTemplatesByTag(ctx context.Context, tag string) ([]models.PromptTemplate, error) {
	return s.repo.GetByTag(ctx, tag)
}

func (s *PromptTemplateService) UpdatePromptTemplate(ctx context.Context, template *models.PromptTemplate) error {
	return s.repo.Update(ctx, template)
}

func (s *PromptTemplateService) DeletePromptTemplate(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.Delete(ctx, objectID)
}

func isValidURL(str string) bool {
	if str == "" {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func fetchJSONData(ctx context.Context, urlStr string) (interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %w", urlStr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from %s: %d", urlStr, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON data: %w", err)
	}

	return data, nil
}

func formatJSONData(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("{error: %s}", err.Error())
	}
	return string(jsonBytes)
}

// FetchFuturesMarketSentiment 从 URL 获取期货数据并返回市场情绪摘要
func (s *PromptTemplateService) FetchFuturesMarketSentiment(ctx context.Context, urlStr string) (string, error) {
	if !isValidURL(urlStr) {
		return "", fmt.Errorf("invalid URL")
	}
	data, err := fetchJSONData(ctx, urlStr)
	if err != nil {
		return "", err
	}
	return futuresMarketSentiment(data), nil
}

// FuturesMover 期货品种行情精简结构（仅暴露给前端用）
type FuturesMover struct {
	Name string  `json:"name"`
	DM   string  `json:"dm"`
	P    float64 `json:"p"`
	Zdf  float64 `json:"zdf"`
}

// FuturesTopMovers 涨跌榜结果：gainers=涨幅前N，losers=跌幅前N，涨跌平计数
type FuturesTopMovers struct {
	Gainers     []FuturesMover `json:"gainers"`
	Losers      []FuturesMover `json:"losers"`
	GainersCnt  int            `json:"gainers_cnt"`
	LosersCnt   int            `json:"losers_cnt"`
	FlatCnt     int            `json:"flat_cnt"`
	Total       int            `json:"total"`
}

// FetchFuturesTopMovers 从 URL 获取期货数据，返回涨/跌幅各前 limit 个品种
func (s *PromptTemplateService) FetchFuturesTopMovers(ctx context.Context, urlStr string, limit int) (*FuturesTopMovers, error) {
	if !isValidURL(urlStr) {
		return nil, fmt.Errorf("invalid URL")
	}
	if limit <= 0 {
		limit = 9
	}
	data, err := fetchJSONData(ctx, urlStr)
	if err != nil {
		return nil, err
	}
	return extractFuturesTopMovers(data, limit), nil
}

func extractFuturesTopMovers(data interface{}, limit int) *FuturesTopMovers {
	result := &FuturesTopMovers{
		Gainers: []FuturesMover{},
		Losers:  []FuturesMover{},
	}
	obj, ok := data.(map[string]interface{})
	if !ok {
		return result
	}
	listRaw, ok := obj["list"]
	if !ok {
		return result
	}
	list, ok := listRaw.([]interface{})
	if !ok {
		return result
	}

	var items []FuturesMover
	for _, raw := range list {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := m["name"].(string)
		dm, _ := m["dm"].(string)
		p, _ := m["p"].(float64)
		zdf, _ := m["zdf"].(float64)
		if name == "" && dm == "" {
			continue
		}
		items = append(items, FuturesMover{Name: name, DM: dm, P: p, Zdf: zdf})
	}
	result.Total = len(items)
	for _, it := range items {
		switch {
		case it.Zdf > 0:
			result.GainersCnt++
		case it.Zdf < 0:
			result.LosersCnt++
		default:
			result.FlatCnt++
		}
	}

	// 按 zdf 降序复制一份取涨幅前 limit
	sortedDesc := make([]FuturesMover, len(items))
	copy(sortedDesc, items)
	sort.SliceStable(sortedDesc, func(i, j int) bool { return sortedDesc[i].Zdf > sortedDesc[j].Zdf })
	if len(sortedDesc) > limit {
		sortedDesc = sortedDesc[:limit]
	}
	result.Gainers = sortedDesc

	// 按 zdf 升序取跌幅前 limit
	sortedAsc := make([]FuturesMover, len(items))
	copy(sortedAsc, items)
	sort.SliceStable(sortedAsc, func(i, j int) bool { return sortedAsc[i].Zdf < sortedAsc[j].Zdf })
	if len(sortedAsc) > limit {
		sortedAsc = sortedAsc[:limit]
	}
	result.Losers = sortedAsc

	return result
}

func futuresMarketSentiment(data interface{}) string {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return formatJSONData(data)
	}
	listRaw, ok := obj["list"]
	if !ok {
		return formatJSONData(data)
	}
	list, ok := listRaw.([]interface{})
	if !ok {
		return formatJSONData(data)
	}

	var up, down, flat int
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		zdf, ok := m["zdf"].(float64)
		if !ok {
			continue
		}
		switch {
		case zdf > 0:
			up++
		case zdf < 0:
			down++
		default:
			flat++
		}
	}
	total := up + down + flat
	var sentiment string
	switch {
	case up > down*2:
		sentiment = "强势多头"
	case down > up*2:
		sentiment = "强势空头"
	case up > down:
		sentiment = "偏多"
	case down > up:
		sentiment = "偏空"
	default:
		sentiment = "多空均衡"
	}
	return fmt.Sprintf("总计%d个合约：上涨%d个，下跌%d个，持平%d个，市场情绪：%s", total, up, down, flat, sentiment)
}

func (s *PromptTemplateService) GeneratePrompt(ctx context.Context, templateID string) (string, error) {
	template, err := s.GetPromptTemplateByID(ctx, templateID)
	if err != nil {
		return "", fmt.Errorf("failed to get prompt template: %w", err)
	}

	sysConfig, err := s.systemConfigService.GetConfig(ctx)
	if err != nil {
		fmt.Printf("Warning: failed to get system config: %v\n", err)
	}

	beijingTime := time.Now().Add(8 * time.Hour)
	openTime := beijingTime.Add(5 * time.Minute)
	prompt := template.Content

	if sysConfig != nil {
		runtimeParameters := cloneParameters(sysConfig.Parameters)
		if runtimeParameters == nil {
			runtimeParameters = map[string]string{}
		}

		// 优先使用租户参数，覆盖全局参数
		if s.tenantConfigService != nil {
			if tenantParams, err := s.tenantConfigService.GetParameters(ctx); err == nil && len(tenantParams) > 0 {
				runtimeParameters = cloneParameters(tenantParams)
			}
		}

		for key, val := range runtimeParameters {
			placeholder := fmt.Sprintf("{{.%s}}", key)
			if !strings.Contains(prompt, placeholder) {
				continue
			}
			if isValidURL(val) {
				data, fetchErr := fetchJSONData(ctx, val)
				if fetchErr != nil {
					prompt = strings.ReplaceAll(prompt, placeholder, fmt.Sprintf("[参数%s数据获取失败: %s]", key, fetchErr.Error()))
				} else {
					prompt = strings.ReplaceAll(prompt, placeholder, formatJSONData(data))
				}
			} else {
				prompt = strings.ReplaceAll(prompt, placeholder, val)
			}
		}
	}

	// 处理 {{.xxxAIResult}} 占位符，注入对应 tab 的最新分析数据
	prompt = s.replaceAIResultPlaceholders(ctx, prompt)

	prompt += fmt.Sprintf("\n- 当前北京时间: %s"+"，开仓时间: %s", beijingTime.Format("2006-01-02 15:04:05"),openTime.Format("2006-01-02 15:04:05"))


	return prompt, nil
}

// replaceAIResultPlaceholders 替换 {{.futuresAIResult}}、{{.optionsAIResult}}、{{.stockAIResult}} 占位符
func (s *PromptTemplateService) replaceAIResultPlaceholders(ctx context.Context, prompt string) string {
	if s.aiResponseRepo == nil {
		return prompt
	}

	tabMap := map[string]string{
		"{{.futuresAIResult}}": models.TabTagFutures,
		"{{.optionsAIResult}}": models.TabTagOptions,
		"{{.stockAIResult}}":   models.TabTagStock,
	}

	for placeholder, tabTag := range tabMap {
		if !strings.Contains(prompt, placeholder) {
			continue
		}
		batchID, err := s.aiResponseRepo.GetLatestCompletedBatchID(ctx, tabTag)
		if err != nil {
			prompt = strings.ReplaceAll(prompt, placeholder, fmt.Sprintf("[%s 暂无分析数据]", tabTag))
			continue
		}
		responses, err := s.aiResponseRepo.GetCompletedByBatchID(ctx, tabTag, batchID)
		if err != nil || len(responses) == 0 {
			prompt = strings.ReplaceAll(prompt, placeholder, fmt.Sprintf("[%s 暂无分析数据]", tabTag))
			continue
		}
		var sb strings.Builder
		for i, r := range responses {
			sb.WriteString(fmt.Sprintf("=== 模型%d（%s）分析结果 ===\n%s\n\n", i+1, r.ModelName, r.Response))
		}
		prompt = strings.ReplaceAll(prompt, placeholder, sb.String())
	}

	return prompt
}
