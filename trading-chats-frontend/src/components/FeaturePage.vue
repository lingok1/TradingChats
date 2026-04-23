<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useTheme } from '../composables/useTheme'

const { mode } = useTheme()

const emit = defineEmits(['switch-to-futures'])

const isDarkMode = computed(() => mode.value === 'dark')
const isVisible = ref(false)
const robotTiltX = ref('0deg')
const robotTiltY = ref('0deg')
const robotGlareX = ref('50%')
const robotGlareY = ref('18%')
const tradeWords = ['buy', 'sell', 'call', 'put'] as const
const floatingWords = ref<string[]>(shuffleTradeWords())
let floatingWordsTimer: number | null = null

const robotStageStyle = computed(() => ({
  '--robot-tilt-x': robotTiltX.value,
  '--robot-tilt-y': robotTiltY.value,
  '--robot-glare-x': robotGlareX.value,
  '--robot-glare-y': robotGlareY.value,
}))

onMounted(() => {
  // 页面加载后显示内容，添加淡入效果
  setTimeout(() => {
    isVisible.value = true
  }, 100)

  floatingWordsTimer = window.setInterval(() => {
    floatingWords.value = shuffleTradeWords()
  }, 2200)
})

onUnmounted(() => {
  if (floatingWordsTimer !== null) {
    window.clearInterval(floatingWordsTimer)
    floatingWordsTimer = null
  }
})

function shuffleTradeWords() {
  const words = [...tradeWords]
  for (let i = words.length - 1; i > 0; i -= 1) {
    const randomIndex = Math.floor(Math.random() * (i + 1))
    ;[words[i], words[randomIndex]] = [words[randomIndex], words[i]]
  }
  return words
}

function handleSwitchToFutures() {
  emit('switch-to-futures')
}

function handleRobotMove(event: MouseEvent) {
  const currentTarget = event.currentTarget
  if (!(currentTarget instanceof HTMLElement)) return

  const rect = currentTarget.getBoundingClientRect()
  const percentX = (event.clientX - rect.left) / rect.width
  const percentY = (event.clientY - rect.top) / rect.height

  robotTiltY.value = `${(percentX - 0.5) * 16}deg`
  robotTiltX.value = `${(0.5 - percentY) * 14}deg`
  robotGlareX.value = `${18 + percentX * 64}%`
  robotGlareY.value = `${8 + percentY * 50}%`
}

function handleRobotLeave() {
  robotTiltX.value = '0deg'
  robotTiltY.value = '0deg'
  robotGlareX.value = '50%'
  robotGlareY.value = '18%'
}

function getTradeWordTone(word: string) {
  return word === 'buy' || word === 'call' ? 'tone-positive' : 'tone-negative'
}
</script>

<template>
  <div class="feature-page" :class="{ 'dark-mode': isDarkMode, visible: isVisible }">
    <!-- 英雄区域 -->
    <section class="hero-section">
      <div class="hero-overlay"></div>
      <div class="hero-flow hero-flow-red"></div>
      <div class="hero-flow hero-flow-green"></div>
      <div class="hero-content">
        <div class="hero-text">
          <h1 class="hero-title">
            <span class="hero-title-primary">AI辅助</span>
            <span class="hero-title-secondary">期货期权交易<span class="profit-text">获得盈利</span></span>
          </h1>
          <p class="hero-subtitle">智能分析市场动态，科学制定交易策略，让投资更专业、更高效</p>
          <div class="hero-buttons">
            <el-button type="primary" size="large" class="cta-button" @click="handleSwitchToFutures">
              开始使用
              <svg
                class="button-icon"
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M5 12h14" />
                <path d="m12 5 7 7-7 7" />
              </svg>
            </el-button>
            <el-button size="large" class="learn-more-button" @click="handleSwitchToFutures">了解更多</el-button>
          </div>
        </div>
        <div class="hero-image">
          <div class="robot-stage" :style="robotStageStyle" @mousemove="handleRobotMove" @mouseleave="handleRobotLeave">
            <div class="robot-stage__glare"></div>
            <div class="robot-stage__core">
              <div class="robot-glow robot-glow-left"></div>
              <div class="robot-glow robot-glow-right"></div>

              <div class="robot-orbit robot-orbit-one"></div>
              <div class="robot-orbit robot-orbit-two"></div>

              <div class="robot-signal robot-signal-left" :class="getTradeWordTone(floatingWords[0])">
                <strong>{{ floatingWords[0] }}</strong>
              </div>
              <div class="robot-signal robot-signal-right" :class="getTradeWordTone(floatingWords[1])">
                <strong>{{ floatingWords[1] }}</strong>
              </div>

              <div class="robot-tag robot-tag-top" :class="getTradeWordTone(floatingWords[2])">{{ floatingWords[2] }}</div>
              <div class="robot-tag robot-tag-bottom" :class="getTradeWordTone(floatingWords[3])">{{ floatingWords[3] }}</div>

              <svg class="robot-figure" viewBox="0 0 320 520" aria-label="人形机器人" role="img">
              <defs>
                <linearGradient id="robotBody" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#f8fbff" />
                  <stop offset="45%" stop-color="#d9e9ff" />
                  <stop offset="100%" stop-color="#7cc4ff" />
                </linearGradient>
                <linearGradient id="robotChrome" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#ffffff" />
                  <stop offset="22%" stop-color="#edf5ff" />
                  <stop offset="54%" stop-color="#b7d8ff" />
                  <stop offset="100%" stop-color="#79bbff" />
                </linearGradient>
                <linearGradient id="robotChromeDark" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#2a5684" />
                  <stop offset="100%" stop-color="#0a1d34" />
                </linearGradient>
                <linearGradient id="robotGlow" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stop-color="#8affea" />
                  <stop offset="100%" stop-color="#4d9fff" />
                </linearGradient>
                <linearGradient id="robotGlass" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#0a1e38" />
                  <stop offset="100%" stop-color="#123961" />
                </linearGradient>
                <linearGradient id="robotShadow" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#102948" />
                  <stop offset="100%" stop-color="#09182d" />
                </linearGradient>
                <radialGradient id="robotCore" cx="50%" cy="40%" r="62%">
                  <stop offset="0%" stop-color="#b0fff2" />
                  <stop offset="45%" stop-color="#66dfff" />
                  <stop offset="100%" stop-color="#19456b" />
                </radialGradient>
                <filter id="robotGlowFilter" x="-60%" y="-60%" width="220%" height="220%">
                  <feDropShadow dx="0" dy="0" stdDeviation="10" flood-color="#6cecff" flood-opacity="0.45" />
                </filter>
              </defs>

              <ellipse cx="160" cy="486" rx="86" ry="18" fill="rgba(87, 173, 255, 0.18)" />
              <ellipse cx="160" cy="488" rx="120" ry="26" fill="rgba(77, 159, 255, 0.08)" />

              <g class="robot-float">
                <path d="M130 64 Q160 42 190 64" fill="none" stroke="url(#robotGlow)" stroke-width="8" stroke-linecap="round" />
                <circle cx="160" cy="52" r="8" fill="#8affea" />

                <ellipse cx="160" cy="122" rx="70" ry="78" fill="rgba(104, 183, 255, 0.08)" filter="url(#robotGlowFilter)" />
                <rect x="104" y="72" width="112" height="126" rx="52" fill="url(#robotChrome)" />
                <path d="M112 98 Q160 70 208 96 V132 Q160 154 112 132 Z" fill="rgba(255,255,255,0.26)" />
                <rect x="96" y="104" width="16" height="46" rx="8" fill="url(#robotChromeDark)" />
                <rect x="208" y="104" width="16" height="46" rx="8" fill="url(#robotChromeDark)" />
                <rect x="116" y="92" width="88" height="60" rx="30" fill="url(#robotGlass)" />
                <path d="M128 103 H193" stroke="rgba(255,255,255,0.2)" stroke-width="4" stroke-linecap="round" />
                <rect x="132" y="112" width="24" height="10" rx="5" fill="#8affea" filter="url(#robotGlowFilter)" />
                <rect x="164" y="112" width="24" height="10" rx="5" fill="#8affea" filter="url(#robotGlowFilter)" />
                <path d="M138 146 Q160 160 182 146" fill="none" stroke="url(#robotGlow)" stroke-width="6" stroke-linecap="round" />
                <path d="M114 168 H206" stroke="rgba(255,255,255,0.24)" stroke-width="6" stroke-linecap="round" />

                <rect x="130" y="206" width="16" height="32" rx="8" fill="url(#robotChrome)" />
                <rect x="174" y="206" width="16" height="32" rx="8" fill="url(#robotChrome)" />
                <path d="M92 236 Q160 208 228 236 L248 314 Q160 364 72 314 Z" fill="url(#robotBody)" />
                <path d="M110 250 Q160 232 210 250 L198 340 Q160 352 122 340 Z" fill="url(#robotChromeDark)" />
                <path d="M94 242 Q160 224 226 242" fill="none" stroke="rgba(255,255,255,0.28)" stroke-width="6" stroke-linecap="round" />
                <rect x="124" y="262" width="72" height="20" rx="10" fill="rgba(255,255,255,0.14)" />
                <path d="M116 298 H204" stroke="url(#robotGlow)" stroke-width="8" stroke-linecap="round" />
                <path d="M130 322 H190" stroke="rgba(215,238,255,0.9)" stroke-width="7" stroke-linecap="round" />
                <circle cx="160" cy="352" r="22" fill="url(#robotCore)" filter="url(#robotGlowFilter)" />

                <circle cx="92" cy="258" r="24" fill="url(#robotChrome)" />
                <circle cx="228" cy="258" r="24" fill="url(#robotChrome)" />
                <rect x="48" y="258" width="38" height="104" rx="19" fill="url(#robotChrome)" />
                <rect x="234" y="258" width="38" height="104" rx="19" fill="url(#robotChrome)" />
                <rect x="38" y="346" width="50" height="84" rx="25" fill="url(#robotChrome)" transform="rotate(16 63 388)" />
                <rect x="232" y="346" width="50" height="84" rx="25" fill="url(#robotChrome)" transform="rotate(-16 257 388)" />
                <rect x="48" y="276" width="14" height="56" rx="7" fill="rgba(255,255,255,0.38)" />
                <rect x="258" y="276" width="14" height="56" rx="7" fill="rgba(255,255,255,0.22)" />
                <circle cx="62" cy="434" r="22" fill="url(#robotChrome)" />
                <circle cx="258" cy="434" r="22" fill="url(#robotChrome)" />

                <path d="M126 372 H194" stroke="rgba(255,255,255,0.2)" stroke-width="6" stroke-linecap="round" />
                <path d="M130 376 L112 452 Q160 474 208 452 L190 376 Z" fill="url(#robotChromeDark)" />
                <rect x="110" y="372" width="38" height="104" rx="19" fill="url(#robotChrome)" />
                <rect x="172" y="372" width="38" height="104" rx="19" fill="url(#robotChrome)" />
                <rect x="116" y="392" width="12" height="62" rx="6" fill="rgba(255,255,255,0.34)" />
                <rect x="188" y="392" width="12" height="62" rx="6" fill="rgba(255,255,255,0.2)" />
                <path d="M106 458 H150" stroke="rgba(255,255,255,0.18)" stroke-width="5" stroke-linecap="round" />
                <path d="M170 458 H214" stroke="rgba(255,255,255,0.18)" stroke-width="5" stroke-linecap="round" />
                <rect x="100" y="458" width="54" height="24" rx="12" fill="url(#robotChrome)" />
                <rect x="166" y="458" width="54" height="24" rx="12" fill="url(#robotChrome)" />

                <path d="M102 112 Q92 130 98 156" fill="none" stroke="rgba(138,255,234,0.45)" stroke-width="5" stroke-linecap="round" />
                <path d="M218 112 Q228 130 222 156" fill="none" stroke="rgba(138,255,234,0.45)" stroke-width="5" stroke-linecap="round" />
                <path d="M82 248 Q58 286 68 338" fill="none" stroke="rgba(77,159,255,0.16)" stroke-width="3" stroke-linecap="round" />
                <path d="M238 248 Q262 286 252 338" fill="none" stroke="rgba(77,159,255,0.16)" stroke-width="3" stroke-linecap="round" />
              </g>
            </svg>

            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 功能介绍 -->
    <section class="features-section">
      <div class="section-header">
        <h2 class="section-title">核心功能</h2>
        <p class="section-subtitle">我们提供全方位的AI辅助交易工具，帮助您在期货期权市场获得优势</p>
      </div>
      <div class="features-grid">
        <!-- 期货期权分析 -->
        <div class="feature-card">
          <div class="feature-icon futures-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
            </svg>
          </div>
          <h3 class="feature-title">期货期权分析</h3>
          <p class="feature-description">
            在日间和夜盘定时分析期货期权品种给出个人投资者建议，期货每5分钟分析90个品种，期权每10分钟分析，供投资者进行参考。
          </p>
        </div>

        <!-- 新闻追踪 -->
        <div class="feature-card">
          <div class="feature-icon news-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="8" x2="12" y2="12" />
              <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
          </div>
          <h3 class="feature-title">新闻追踪</h3>
          <p class="feature-description">实时追踪期货期权市场动态，把握投资先机。</p>
        </div>

        <!-- 交易计划 -->
        <div class="feature-card">
          <div class="feature-icon plan-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
              <line x1="16" y1="2" x2="16" y2="6" />
              <line x1="8" y1="2" x2="8" y2="6" />
              <line x1="3" y1="10" x2="21" y2="10" />
            </svg>
          </div>
          <h3 class="feature-title">交易计划</h3>
          <p class="feature-description">
            计划你的交易，交易你的计划，制定期货期权交易计划，克服人性弱点，建立交易的一致性，提高决策效率，减轻心理压力。
          </p>
        </div>

        <!-- 持仓分析 -->
        <div class="feature-card">
          <div class="feature-icon position-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
            </svg>
          </div>
          <h3 class="feature-title">持仓分析</h3>
          <p class="feature-description">AI给出持仓分析和建议，可通过图片导入和手动添加期货期权持仓。</p>
        </div>
      </div>
    </section>

    <!-- 优势介绍 -->
    <section class="advantages-section">
      <div class="section-header">
        <h2 class="section-title">我们的优势</h2>
        <p class="section-subtitle">专业的AI技术，为您的投资决策提供强大支持</p>
      </div>
      <div class="advantages-content">
        <div class="advantage-item">
          <div class="advantage-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="32"
              height="32"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
              <polyline points="22 4 12 14.01 9 11.01" />
            </svg>
          </div>
          <h3 class="advantage-title">智能分析</h3>
          <p class="advantage-description">利用AI技术对市场数据进行深度分析，提供专业的投资建议。</p>
        </div>
        <div class="advantage-item">
          <div class="advantage-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="32"
              height="32"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="10" />
              <polyline points="12 6 12 12 16 14" />
            </svg>
          </div>
          <h3 class="advantage-title">实时更新</h3>
          <p class="advantage-description">定时更新分析结果，确保你获取最新的市场动态和投资机会。</p>
        </div>
        <div class="advantage-item">
          <div class="advantage-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="32"
              height="32"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
              <path d="M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
          </div>
          <h3 class="advantage-title">专业指导</h3>
          <p class="advantage-description">基于专业的交易策略和风险管理，帮助你制定科学的投资计划。</p>
        </div>
      </div>
    </section>

    <!-- 数据展示 -->
    <section class="data-section">
      <div class="data-content">
        <div class="data-item">
          <div class="data-number">80+</div>
          <div class="data-label">期货品种</div>
        </div>
        <div class="data-item">
          <div class="data-number">5分钟</div>
          <div class="data-label">分析频率</div>
        </div>
        <div class="data-item">
          <div class="data-number">交易日</div>
          <div class="data-label">全天候监控</div>
        </div>
        <div class="data-item">
          <div class="data-number">AI辅助</div>
          <div class="data-label">智能分析</div>
        </div>
      </div>
    </section>

    <!-- 行动号召 -->
    <section class="cta-section">
      <div class="cta-content">
        <h2 class="cta-title">开始你的智能投资之旅</h2>
        <p class="cta-description">加入我们，体验AI辅助交易的优势，让投资决策更加科学、高效。</p>
        <el-button type="primary" size="large" class="cta-button" @click="handleSwitchToFutures">
          立即开始
          <svg
            class="button-icon"
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M5 12h14" />
            <path d="m12 5 7 7-7 7" />
          </svg>
        </el-button>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* 基础样式 */
.feature-page {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  color: var(--el-text-color-primary);
  opacity: 0;
  transition: opacity 0.8s ease;
}

.feature-page.visible {
  opacity: 1;
}

/* 英雄区域 */
.hero-section {
  position: relative;
  padding: 80px 0;
  background: linear-gradient(135deg, #0a192f 0%, #172a45 100%);
  color: white;
  overflow: hidden;
}

.hero-flow {
  position: absolute;
  z-index: 0;
  width: min(78vw, 980px);
  height: 360px;
  border-top: 10px solid transparent;
  border-radius: 50%;
  opacity: 0.86;
  pointer-events: none;
  filter: blur(0.2px);
}

.hero-flow::before {
  content: '';
  position: absolute;
  inset: -6px -6px auto;
  height: 10px;
  border-radius: 999px;
  background: linear-gradient(90deg, transparent, currentColor 22%, currentColor 78%, transparent);
  opacity: 0.52;
  filter: blur(5px);
}

.hero-flow::after {
  content: '';
  position: absolute;
  inset: -6px -14px auto;
  height: 22px;
  border-radius: 999px;
  background: currentColor;
  filter: blur(20px);
  opacity: 0.78;
}

.hero-flow-red {
  top: 36px;
  left: -210px;
  color: rgba(255, 86, 108, 0.92);
  border-top-color: rgba(255, 86, 108, 0.94);
  transform: rotate(-9deg);
  animation: heroFlowRed 12.5s ease-in-out infinite alternate;
}

.hero-flow-green {
  right: -210px;
  bottom: 16px;
  color: rgba(83, 236, 160, 0.9);
  border-top-color: rgba(83, 236, 160, 0.92);
  transform: rotate(11deg);
  animation: heroFlowGreen 13.5s ease-in-out infinite alternate;
}

.hero-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxkZWZzPjxwYXR0ZXJuIGlkPSJwYXR0ZXJuIiB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHBhdHRlcm5Vbml0cz0idXNlclNwYWNlT25Vc2UiIHBhdHRlcm5UcmFuc2Zvcm09InJvdGF0ZSgxNSkiPjxwYXRoIGQ9Ik0gNDAgMCBMIDAgMCAwIDQwIiBmaWxsPSJub25lIiBzdHJva2U9IiMzYjQxNTYiIHN0cm9rZS13aWR0aD0iMC41Ii8+PC9wYXR0ZXJuPjwvZGVmcz48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSJ1cmwoI3BhdHRlcm4pIiAvPjwvc3ZnPg==');
  opacity: 0.3;
}

.hero-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 40px;
  position: relative;
  z-index: 1;
}

.hero-text {
  flex: 1;
  max-width: 600px;
}

.hero-title {
  font-size: 40px;
  font-weight: 700;
  margin-bottom: 16px;
  line-height: 1.2;
  letter-spacing: -0.5px;
}

.hero-title-primary {
  color: #64ffda;
  display: block;
  margin-bottom: 4px;
}

.hero-title-secondary {
  display: block;
}

.profit-text {
  color: #ff6b6b;
  margin-left: 8px;
}

.hero-subtitle {
  font-size: 18px;
  margin-bottom: 30px;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.4;
}

.hero-buttons {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.cta-button {
  background: transparent;
  border: 1px solid #64ffda;
  color: #64ffda;
  padding: 12px 24px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.cta-button:hover {
  background: rgba(100, 255, 218, 0.1);
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(100, 255, 218, 0.2);
}

.learn-more-button {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: white;
  padding: 12px 24px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.learn-more-button:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateY(-2px);
  border-color: rgba(255, 255, 255, 0.4);
}

.button-icon {
  transition: transform 0.3s ease;
}

.cta-button:hover .button-icon {
  transform: translateX(4px);
}

.hero-image {
  flex: 1;
  max-width: 450px;
}

.robot-stage {
  --robot-tilt-x: 0deg;
  --robot-tilt-y: 0deg;
  --robot-glare-x: 50%;
  --robot-glare-y: 18%;
  position: relative;
  min-height: 520px;
  border-radius: 24px;
  overflow: hidden;
  border: 1px solid rgba(100, 255, 218, 0.18);
  background:
    radial-gradient(circle at 50% 18%, rgba(123, 220, 255, 0.14), transparent 26%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.02)),
    rgba(9, 22, 39, 0.78);
  box-shadow:
    0 24px 50px rgba(0, 0, 0, 0.34),
    inset 0 1px 0 rgba(255, 255, 255, 0.08);
  transform-style: preserve-3d;
  transition: transform 0.25s ease;
}

.robot-stage::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(100, 255, 218, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(100, 255, 218, 0.08) 1px, transparent 1px);
  background-size: 34px 34px;
  opacity: 0.28;
}

.robot-stage::after {
  content: '';
  position: absolute;
  inset: -30% 0 auto;
  height: 42%;
  background: linear-gradient(180deg, rgba(100, 255, 218, 0.12), transparent);
  animation: scanLine 6.2s linear infinite;
}

.robot-stage__core {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 520px;
  transform-style: preserve-3d;
  transform: perspective(1200px) rotateX(var(--robot-tilt-x)) rotateY(var(--robot-tilt-y));
  transition: transform 0.22s ease;
}

.robot-stage__glare {
  position: absolute;
  inset: 0;
  z-index: 3;
  pointer-events: none;
  background:
    radial-gradient(circle at var(--robot-glare-x) var(--robot-glare-y), rgba(255, 255, 255, 0.28), transparent 18%),
    linear-gradient(135deg, rgba(255, 255, 255, 0.02), transparent 42%);
  mix-blend-mode: screen;
}

.robot-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(42px);
  pointer-events: none;
}

.robot-glow-left {
  width: 180px;
  height: 180px;
  left: 26px;
  top: 112px;
  background: rgba(100, 255, 218, 0.2);
}

.robot-glow-right {
  width: 160px;
  height: 160px;
  right: 32px;
  top: 84px;
  background: rgba(77, 159, 255, 0.2);
}

.robot-orbit {
  position: absolute;
  left: 50%;
  top: 48%;
  border: 1px solid rgba(100, 255, 218, 0.16);
  border-radius: 50%;
  transform: translate(-50%, -50%);
}

.robot-orbit-one {
  width: 240px;
  height: 360px;
  animation: orbitSpin 16s linear infinite;
}

.robot-orbit-two {
  width: 300px;
  height: 430px;
  border-color: rgba(77, 159, 255, 0.14);
  animation: orbitSpin 24s linear infinite reverse;
}

.robot-signal {
  position: absolute;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid rgba(100, 255, 218, 0.16);
  background: rgba(7, 18, 32, 0.74);
  backdrop-filter: blur(10px);
  box-shadow: 0 16px 30px rgba(0, 0, 0, 0.18);
  transform: translateZ(34px);
}

.robot-signal-left {
  left: 18px;
  top: 28px;
  animation: floatCard 4.6s ease-in-out infinite;
}

.robot-signal-right {
  right: 18px;
  top: 86px;
  animation: floatCard 5.2s ease-in-out infinite reverse;
}

.robot-tag {
  position: absolute;
  z-index: 2;
  padding: 8px 12px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(6, 16, 30, 0.7);
  color: rgba(255, 255, 255, 0.82);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  transform: translateZ(26px);
}

.robot-tag-top {
  left: 50%;
  top: 18px;
  transform: translateX(-50%) translateZ(30px);
}

.robot-tag-bottom {
  left: 20px;
  bottom: 18px;
}

.robot-signal strong {
  font-size: 19px;
  line-height: 1;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.tone-positive {
  border-color: rgba(83, 236, 160, 0.28);
  background: rgba(11, 28, 24, 0.7);
  color: #73ffbc;
  box-shadow: 0 0 22px rgba(83, 236, 160, 0.12);
}

.tone-negative {
  border-color: rgba(255, 86, 108, 0.26);
  background: rgba(35, 15, 22, 0.72);
  color: #ff8a9b;
  box-shadow: 0 0 22px rgba(255, 86, 108, 0.12);
}

.robot-figure {
  position: absolute;
  inset: 24px 20px 8px;
  width: calc(100% - 40px);
  height: calc(100% - 32px);
  z-index: 1;
  transform: translateZ(18px);
  filter: drop-shadow(0 24px 34px rgba(7, 18, 31, 0.38));
}

.robot-float {
  transform-origin: center 54%;
  animation: robotFloat 4.8s ease-in-out infinite;
}



/* 功能介绍 */
.features-section {
  padding: 80px 0;
  background: var(--el-bg-color);
}

.section-header {
  text-align: center;
  margin-bottom: 50px;
}

.section-title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 12px;
  color: var(--el-text-color-primary);
  letter-spacing: -0.5px;
}

.section-subtitle {
  font-size: 16px;
  color: var(--el-text-color-secondary);
  max-width: 600px;
  margin: 0 auto;
  line-height: 1.4;
}

.features-grid {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.feature-card {
  background: var(--el-bg-color);
  border-radius: 8px;
  padding: 30px 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  border: 1px solid var(--el-border-color-light);
  position: relative;
  overflow: hidden;
}

.feature-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: linear-gradient(180deg, #64ffda 0%, #4fd1c5 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.12);
}

.feature-card:hover::before {
  opacity: 1;
}

.feature-icon {
  width: 50px;
  height: 50px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  font-size: 20px;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  color: #3182ce;
  transition: all 0.3s ease;
}

.feature-card:hover .feature-icon {
  transform: scale(1.1);
  background: linear-gradient(135deg, #3182ce 0%, #2c5282 100%);
  color: white;
}

.feature-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--el-text-color-primary);
  transition: color 0.3s ease;
}

.feature-card:hover .feature-title {
  color: #3182ce;
}

.feature-description {
  font-size: 14px;
  line-height: 1.5;
  color: var(--el-text-color-secondary);
}

/* 优势介绍 */
.advantages-section {
  padding: 80px 0;
  background: linear-gradient(135deg, #f8fafc 0%, #edf2f7 100%);
}

.advantages-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 30px;
}

.advantage-item {
  text-align: center;
  padding: 30px 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  border: 1px solid var(--el-border-color-light);
}

.advantage-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.12);
}

.advantage-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #3182ce 0%, #2c5282 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  font-size: 24px;
  color: white;
  transition: all 0.3s ease;
}

.advantage-item:hover .advantage-icon {
  transform: scale(1.1);
  background: linear-gradient(135deg, #64ffda 0%, #4fd1c5 100%);
  color: #0a192f;
}

.advantage-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--el-text-color-primary);
}

.advantage-description {
  font-size: 13px;
  line-height: 1.5;
  color: var(--el-text-color-secondary);
}

/* 数据展示 */
.data-section {
  padding: 60px 0;
  background: var(--el-bg-color);
  border-top: 1px solid var(--el-border-color-light);
  border-bottom: 1px solid var(--el-border-color-light);
}

.data-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  justify-content: space-around;
  flex-wrap: wrap;
  gap: 30px;
}

.data-item {
  text-align: center;
  flex: 1;
  min-width: 120px;
}

.data-number {
  font-size: 36px;
  font-weight: 700;
  color: #3182ce;
  margin-bottom: 4px;
  line-height: 1;
}

.data-label {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  text-transform: uppercase;
  letter-spacing: 1px;
}

/* 行动号召 */
.cta-section {
  padding: 80px 0;
  background: linear-gradient(135deg, #0a192f 0%, #172a45 100%);
  color: white;
  text-align: center;
}

.cta-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 20px;
}

.cta-title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 16px;
  letter-spacing: -0.5px;
}

.cta-description {
  font-size: 16px;
  margin-bottom: 30px;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.4;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .hero-flow {
    width: min(118vw, 960px);
    height: 340px;
  }

  .hero-content {
    flex-direction: column;
    text-align: center;
  }

  .hero-text {
    max-width: 100%;
  }

  .hero-image {
    max-width: 100%;
    margin-top: 30px;
  }

  .robot-stage {
    width: min(100%, 520px);
    margin: 0 auto;
  }

  .robot-stage__core {
    transform: none;
  }

  .hero-title {
    font-size: 32px;
  }
}

@media (max-width: 768px) {
  .hero-section {
    padding: 44px 0 60px;
  }

  .hero-flow-red {
    width: 152vw;
    height: 290px;
    left: -42vw;
    top: 68px;
  }

  .hero-flow-green {
    width: 156vw;
    height: 290px;
    right: -44vw;
    bottom: 26px;
  }

  .hero-title {
    font-size: 24px;
  }

  .hero-subtitle {
    font-size: 14px;
  }

  .hero-content {
    gap: 20px;
  }

  .hero-image {
    order: -1;
    width: 100%;
    margin-top: 0;
  }

  .hero-text {
    order: 1;
  }

  .hero-buttons {
    flex-direction: column;
    align-items: center;
  }

  .robot-stage {
    min-height: 340px;
    border-radius: 18px;
  }

  .robot-stage__core {
    min-height: 340px;
  }

  .robot-figure {
    inset: 10px 8px 2px;
    width: calc(100% - 16px);
    height: calc(100% - 12px);
  }

  .robot-signal {
    padding: 8px 10px;
    border-radius: 10px;
  }

  .robot-signal strong {
    font-size: 14px;
    letter-spacing: 0.1em;
  }

  .robot-signal-left {
    left: 8px;
    top: 10px;
  }

  .robot-signal-right {
    right: 8px;
    top: 56px;
  }

  .robot-tag {
    padding: 6px 8px;
    font-size: 10px;
    letter-spacing: 0.1em;
  }

  .robot-tag-top {
    top: 10px;
  }

  .robot-tag-bottom {
    left: 8px;
    bottom: 8px;
  }

  .features-section,
  .advantages-section,
  .data-section,
  .cta-section {
    padding: 60px 0;
  }

  .section-title {
    font-size: 24px;
  }

  .section-subtitle {
    font-size: 14px;
  }

  .features-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .feature-card {
    padding: 24px 16px;
  }

  .advantages-content {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .advantage-item {
    padding: 24px 16px;
  }

  .data-content {
    flex-direction: column;
    align-items: center;
  }

  .data-item {
    min-width: 180px;
  }

  .cta-title {
    font-size: 24px;
  }

  .cta-description {
    font-size: 14px;
  }
}

@keyframes robotFloat {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-12px);
  }
}

@keyframes orbitSpin {
  from {
    transform: translate(-50%, -50%) rotate(0deg);
  }
  to {
    transform: translate(-50%, -50%) rotate(360deg);
  }
}

@keyframes floatCard {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-8px);
  }
}

@keyframes scanLine {
  0% {
    transform: translateY(-120%);
  }
  100% {
    transform: translateY(280%);
  }
}

@keyframes panelBarPulse {
  0%,
  100% {
    opacity: 0.85;
    transform: scaleY(0.94);
  }
  50% {
    opacity: 1;
    transform: scaleY(1.04);
  }
}

@keyframes heroFlowRed {
  0% {
    transform: translate3d(0, 0, 0) rotate(-9deg) scaleX(1);
  }
  100% {
    transform: translate3d(90px, 26px, 0) rotate(-5deg) scaleX(1.05);
  }
}

@keyframes heroFlowGreen {
  0% {
    transform: translate3d(0, 0, 0) rotate(11deg) scaleX(1);
  }
  100% {
    transform: translate3d(-104px, -30px, 0) rotate(7deg) scaleX(1.06);
  }
}

/* 深色模式适配 */
.dark-mode .features-section {
  background: #1a202c;
}

.dark-mode .feature-card {
  background: #2d3748;
  border-color: #4a5568;
}

.dark-mode .feature-icon {
  background: linear-gradient(135deg, #2d3748 0%, #4a5568 100%);
  color: #64ffda;
}

.dark-mode .feature-card:hover .feature-icon {
  background: linear-gradient(135deg, #64ffda 0%, #4fd1c5 100%);
  color: #0a192f;
}

.dark-mode .advantages-section {
  background: #1a202c;
}

.dark-mode .advantage-item {
  background: #2d3748;
  border-color: #4a5568;
}

.dark-mode .data-section {
  background: #1a202c;
  border-color: #4a5568;
}
</style>
