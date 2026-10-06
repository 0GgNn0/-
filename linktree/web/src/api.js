async function request(url, options = {}) {
  const res = await fetch(url, {
    headers: { 'Content-Type': 'application/json', ...(options.headers || {}) },
    ...options
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    const err = new Error(data.message || 'HTTP ' + res.status)
    err.status = res.status
    throw err
  }
  return data
}

export const api = {
  profile: () => request('/api/public/profile'),
  links: () => request('/api/public/links'),
  resume: () => request('/api/public/resume'),
  works: () => request('/api/works/public'),
  categories: () => request('/api/categories'),
  experiences: () => request('/api/experiences')
}

// SSE 流式问答：data: {"type":"start|delta|done|error"}
export async function askAI(question, onDelta) {
  const res = await fetch('/api/ask', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ question })
  })
  if (!res.ok) {
    const data = await res.json().catch(() => ({}))
    const err = new Error(data.message || '请求失败 (' + res.status + ')')
    err.status = res.status
    throw err
  }
  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  let buf = ''
  let full = ''
  let errorMsg = ''
  while (true) {
    const { done, value } = await reader.read()
    if (done) break
    buf += decoder.decode(value, { stream: true })
    let idx
    while ((idx = buf.indexOf('\n\n')) !== -1) {
      const raw = buf.slice(0, idx)
      buf = buf.slice(idx + 2)
      if (!raw.startsWith('data: ')) continue
      const payload = raw.slice(6).trim()
      if (!payload || payload === '[DONE]') continue
      try {
        const evt = JSON.parse(payload)
        if (evt.type === 'delta' && evt.content) {
          full += evt.content
          onDelta(full)
        } else if (evt.type === 'error') {
          errorMsg = evt.message || 'AI 服务出错'
        }
      } catch (_) { /* skip malformed */ }
    }
  }
  if (errorMsg) {
    const err = new Error(errorMsg)
    err.status = 429
    throw err
  }
  return full
}
