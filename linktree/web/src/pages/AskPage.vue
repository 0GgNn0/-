<script setup>
import { ref, nextTick, onMounted } from 'vue'
import { NButton, NInput, NSpin, useMessage } from 'naive-ui'
import { askAI } from '../api'
import { profile } from '../profile'

const message = useMessage()

const question = ref('')
const loading = ref(false)
const messages = ref([]) // { role: 'user'|'assistant', content }
const scrollBox = ref(null)

const suggestions = [
  '介绍一下你的技术栈',
  '你做过哪些项目？',
  '你最有代表性的作品是什么？',
  '如何联系你？'
]

function scrollToBottom() {
  nextTick(() => {
    if (scrollBox.value) scrollBox.value.scrollTop = scrollBox.value.scrollHeight
  })
}

async function send(q) {
  const text = (q || question.value).trim()
  if (!text || loading.value) return
  question.value = ''
  messages.value.push({ role: 'user', content: text })
  messages.value.push({ role: 'assistant', content: '' })
  loading.value = true
  scrollToBottom()

  try {
    await askAI(text, (full) => {
      messages.value[messages.value.length - 1].content = full
      scrollToBottom()
    })
    if (!messages.value[messages.value.length - 1].content) {
      messages.value[messages.value.length - 1].content = '（无回复）'
    }
  } catch (err) {
    messages.value[messages.value.length - 1].content = ''
    message.error(err.message || '请求失败，请稍后再试')
  } finally {
    loading.value = false
    scrollToBottom()
  }
}

function onKeydown(e) {
  if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
    e.preventDefault()
    send()
  }
}

onMounted(scrollToBottom)
</script>

<template>
  <div class="ask-page">
    <div class="reveal" style="--d:0s">
      <router-link to="/" class="page-back">← 返回首页</router-link>
    </div>

    <div class="chat glass reveal" style="--d:.1s">
      <div class="kicker mono">AI · Ask About Me</div>
      <h1 class="h1">问 AI 关于我</h1>

      <!-- 空状态：欢迎 + 推荐问题 -->
      <div v-if="!messages.length" class="welcome">
        <p class="welcome-text">
          你好，我是 {{ profile.display_name }} 的 AI 助手 👋<br>
          关于他的经历、技能和作品，尽管问我。
        </p>
        <div class="suggestions">
          <button v-for="s in suggestions" :key="s" class="sug-chip" @click="send(s)">{{ s }}</button>
        </div>
      </div>

      <!-- 消息列表 -->
      <div v-else ref="scrollBox" class="messages">
        <div v-for="(m, i) in messages" :key="i" class="msg" :class="m.role">
          <div class="bubble" :class="{ typing: m.role === 'assistant' && !m.content && loading && i === messages.length - 1 }">
            <template v-if="m.role === 'assistant' && !m.content && loading && i === messages.length - 1">
              <span class="dot"></span><span class="dot"></span><span class="dot"></span>
            </template>
            <span v-else style="white-space: pre-wrap">{{ m.content }}</span>
          </div>
        </div>
        <div v-if="loading" class="scroll-anchor"></div>
      </div>

      <!-- 输入区 -->
      <div class="input-bar">
        <n-input
          v-model:value="question"
          type="textarea"
          :rows="1"
          :autosize="{ minRows: 1, maxRows: 5 }"
          placeholder="问我任何问题…（Enter 发送）"
          :disabled="loading"
          @keydown="onKeydown"
        />
        <n-button
          type="primary"
          :loading="loading"
          :disabled="!question.trim()"
          @click="send()"
        >
          发送
        </n-button>
      </div>
      <div class="hint mono">AI 回答由大模型生成，仅供参考 · 每分钟 3 次</div>
    </div>
  </div>
</template>

<style scoped>
.ask-page { max-width: 680px; margin: 0 auto; padding: 28px 24px 72px; }

.page-back {
  display: inline-block;
  font-size: 13px; color: rgba(90, 90, 100, 0.95);
  text-decoration: none; padding: 6px 14px;
  border-radius: 999px; border: 1px solid rgba(0, 0, 0, 0.1);
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .page-back { background: rgba(24, 24, 27, 0.6); color: rgba(190, 190, 200, 0.95); border-color: rgba(255, 255, 255, 0.12); }
.page-back:hover { color: #10b981; border-color: #10b981; transform: translateY(-1px); }

.chat { margin-top: 20px; border-radius: 16px; padding: 32px; display: flex; flex-direction: column; min-height: 480px; }
.kicker { font-size: 11px; letter-spacing: 0.14em; text-transform: uppercase; color: #10b981; margin-bottom: 10px; }
.h1 { font-size: clamp(28px, 5vw, 40px); font-weight: 800; letter-spacing: -0.03em; line-height: 1.05; }

.welcome { margin-top: 28px; }
.welcome-text { font-size: 15px; line-height: 1.8; color: rgba(90, 90, 100, 0.95); }
body[data-theme="dark"] .welcome-text { color: rgba(180, 180, 190, 0.95); }
.suggestions { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 18px; }
.sug-chip {
  padding: 8px 16px;
  border-radius: 999px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  font-size: 13px;
  color: rgba(90, 90, 100, 0.95);
  cursor: pointer;
  font-family: inherit;
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
body[data-theme="dark"] .sug-chip { background: rgba(24, 24, 27, 0.6); color: rgba(190, 190, 200, 0.95); border-color: rgba(255, 255, 255, 0.12); }
.sug-chip:hover { border-color: #10b981; color: #10b981; transform: translateY(-1px); }

.messages { margin-top: 24px; flex: 1; max-height: 420px; overflow-y: auto; display: flex; flex-direction: column; gap: 12px; padding-right: 4px; }
.msg { display: flex; }
.msg.user { justify-content: flex-end; }
.msg.assistant { justify-content: flex-start; }
.bubble {
  max-width: 82%;
  padding: 10px 16px;
  border-radius: 16px;
  font-size: 14px;
  line-height: 1.7;
}
.msg.user .bubble {
  background: #10b981;
  color: #fff;
  border-bottom-right-radius: 4px;
}
.msg.assistant .bubble {
  background: rgba(0, 0, 0, 0.05);
  color: inherit;
  border-bottom-left-radius: 4px;
}
body[data-theme="dark"] .msg.assistant .bubble { background: rgba(255, 255, 255, 0.08); }
.bubble.typing { display: inline-flex; gap: 5px; align-items: center; padding: 14px 16px; }
.dot {
  width: 7px; height: 7px;
  border-radius: 50%;
  background: rgba(120, 120, 130, 0.7);
  animation: blink 1.2s infinite both;
}
.dot:nth-child(2) { animation-delay: 0.2s; }
.dot:nth-child(3) { animation-delay: 0.4s; }
@keyframes blink { 0%, 80%, 100% { opacity: 0.25; } 40% { opacity: 1; } }
.scroll-anchor { height: 1px; }

.input-bar { display: flex; gap: 10px; margin-top: 20px; align-items: flex-end; }
.input-bar :deep(.n-input) { border-radius: 12px; }
.hint { font-size: 11px; color: rgba(120, 120, 130, 0.75); margin-top: 10px; letter-spacing: 0.04em; text-align: center; }

@media (max-width: 768px) {
  .chat { padding: 24px 18px; }
  .bubble { max-width: 90%; }
}
</style>
