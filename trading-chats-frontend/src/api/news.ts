import type { NewsItem, NewsCategory } from './types'

// 模拟新闻分类数据
const mockCategories: NewsCategory[] = [
  { id: '1', name: '宏观经济', code: 'macro', icon: '📊' },
  { id: '2', name: '期货市场', code: 'futures', icon: '📈' },
  { id: '3', name: '期权市场', code: 'options', icon: '⚖️' },
  { id: '4', name: '国际市场', code: 'international', icon: '🌍' },
  { id: '5', name: '政策法规', code: 'policy', icon: '📋' },
  { id: '6', name: '行业动态', code: 'industry', icon: '🏭' }
]

// 模拟新闻数据
const mockNews: NewsItem[] = [
  {
    id: '1',
    title: '央行发布2026年第二季度货币政策执行报告',
    summary: '报告显示，当前经济运行总体平稳，物价水平温和上涨，就业形势稳定。央行将继续实施稳健的货币政策，保持流动性合理充裕。',
    content: '中国人民银行今日发布2026年第二季度货币政策执行报告，报告指出，当前我国经济运行总体平稳，物价水平温和上涨，就业形势稳定。\n\n报告强调，央行将继续实施稳健的货币政策，保持流动性合理充裕，引导市场利率下行，支持实体经济发展。同时，将加强宏观审慎管理，防范系统性金融风险。\n\n对于未来经济走势，报告认为，随着各项稳增长政策措施的落地见效，经济有望保持稳定增长态势。',
    category: '宏观经济',
    source: '央行网站',
    author: '金融市场司',
    publish_time: '2026-04-15 10:00:00',
    image_url: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=central%20bank%20monetary%20policy%20report%20financial%20data%20chart&image_size=landscape_16_9',
    read_count: 12580
  },
  {
    id: '2',
    title: '铜价创近期新高，全球供应紧张持续',
    summary: '受全球经济复苏和供应短缺影响，国际铜价持续上涨，创下近期新高。分析师预计，铜价有望继续保持强势。',
    content: '国际铜价近期持续上涨，创下近6个月新高。受全球经济复苏和供应短缺影响，铜价走势强劲。\n\n据伦敦金属交易所数据显示，LME铜价已突破9000美元/吨关口，较年初上涨超过15%。市场分析认为，全球铜供应紧张状况短期内难以缓解，加之电动汽车和可再生能源领域对铜的需求持续增长，铜价有望继续保持强势。\n\n国内市场方面，沪铜期货价格也跟随上涨，主力合约价格突破70000元/吨。业内人士建议投资者关注相关上市公司的投资机会。',
    category: '期货市场',
    source: '期货日报',
    author: '金属市场分析师',
    publish_time: '2026-04-14 16:30:00',
    image_url: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=copper%20price%20rising%20financial%20chart%20trading%20screen&image_size=landscape_16_9',
    read_count: 8920
  },
  {
    id: '3',
    title: '沪深300期权成交量创新高，市场活跃度提升',
    summary: '随着市场波动加剧，沪深300期权成交量创新高，投资者利用期权工具进行风险管理的需求增加。',
    content: '近期，沪深300期权市场活跃度显著提升，成交量创下历史新高。数据显示，本月沪深300期权日均成交量达到150万手，较上月增长30%。\n\n分析人士认为，市场波动加剧是期权成交量增加的主要原因。投资者通过期权工具进行风险管理和方向性交易的需求明显增加。同时，随着期权市场的不断发展，更多机构投资者开始参与期权交易，进一步提升了市场流动性。\n\n专家建议，投资者在参与期权交易时应充分了解相关风险，合理控制仓位，避免过度投机。',
    category: '期权市场',
    source: '证券时报',
    author: '衍生品市场研究员',
    publish_time: '2026-04-13 09:15:00',
    image_url: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=options%20trading%20volume%20stock%20market%20financial%20data&image_size=landscape_16_9',
    read_count: 6750
  },
  {
    id: '4',
    title: '美联储维持利率不变，市场预期年内将开始降息',
    summary: '美联储在最新货币政策会议上决定维持利率不变，但暗示年内可能开始降息。市场对此反应积极，美股三大指数均上涨。',
    content: '美国联邦储备委员会在最新货币政策会议上决定维持当前利率水平不变，但暗示如果通胀持续回落，年内可能开始降息。\n\n美联储主席鲍威尔在会后新闻发布会上表示，虽然通胀已经有所回落，但仍高于2%的目标水平，因此需要保持谨慎。不过，他也表示，委员会已经开始讨论降息的时机和幅度。\n\n市场对这一消息反应积极，美股三大指数均出现上涨。美元指数下跌，黄金价格上涨。分析人士认为，美联储的表态为市场提供了明确的政策预期，有助于稳定市场情绪。',
    category: '国际市场',
    source: '华尔街日报',
    author: '金融市场记者',
    publish_time: '2026-04-12 22:00:00',
    image_url: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fed%20monetary%20policy%20meeting%20press%20conference&image_size=landscape_16_9',
    read_count: 15320
  },
  {
    id: '5',
    title: '证监会发布《期货和衍生品法实施细则》',
    summary: '证监会近日发布《期货和衍生品法实施细则》，进一步规范期货市场秩序，保护投资者合法权益。',
    content: '中国证监会近日发布《期货和衍生品法实施细则》，自2026年5月1日起施行。细则对期货交易、结算、风险控制等方面做出了详细规定，进一步规范期货市场秩序，保护投资者合法权益。\n\n细则明确了期货公司的业务范围和监管要求，强化了对市场操纵、内幕交易等违法行为的打击力度。同时，细则还完善了期货市场的风险控制机制，提高了市场的抗风险能力。\n\n业内人士认为，细则的发布有利于期货市场的长期健康发展，将为市场参与者提供更加明确的规则指引。',
    category: '政策法规',
    source: '证监会网站',
    author: '法规部',
    publish_time: '2026-04-11 14:00:00',
    image_url: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=securities%20regulatory%20document%20financial%20law%20implementation&image_size=landscape_16_9',
    read_count: 9870
  },
  {
    id: '6',
    title: '新能源汽车销量持续增长，锂资源需求旺盛',
    summary: '2026年第一季度，国内新能源汽车销量同比增长45%，带动锂资源需求持续旺盛，锂价保持高位运行。',
    content: '据中国汽车工业协会数据显示，2026年第一季度，国内新能源汽车销量达到180万辆，同比增长45%，市场渗透率超过35%。\n\n新能源汽车销量的持续增长带动了锂资源需求的旺盛。数据显示，国内锂盐产量虽然有所增加，但仍无法满足快速增长的需求，锂价保持高位运行。\n\n业内分析认为，随着全球新能源汽车产业的快速发展，锂资源将在未来几年持续处于紧平衡状态。相关企业应加强锂资源的勘探和开发，提高供应能力。',
    category: '行业动态',
    source: '中国汽车报',
    author: '产业分析师',
    publish_time: '2026-04-10 11:30:00',
    image_url: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=new%20energy%20vehicles%20lithium%20batteries%20automotive%20industry&image_size=landscape_16_9',
    read_count: 7640
  }
]

// 模拟API延迟
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

export async function getNewsList(category?: string, page: number = 1, pageSize: number = 20) {
  await delay(500) // 模拟网络延迟
  let filteredNews = [...mockNews]
  
  // 按分类筛选
  if (category && category !== '') {
    const categoryName = mockCategories.find(c => c.code === category)?.name
    if (categoryName) {
      filteredNews = filteredNews.filter(news => news.category === categoryName)
    }
  }
  
  // 分页处理
  const start = (page - 1) * pageSize
  const end = start + pageSize
  return filteredNews.slice(start, end)
}

export async function getNewsDetail(id: string) {
  await delay(300) // 模拟网络延迟
  const news = mockNews.find(item => item.id === id)
  if (!news) {
    throw new Error('新闻不存在')
  }
  return news
}

export async function getNewsCategories() {
  await delay(200) // 模拟网络延迟
  return mockCategories
}
