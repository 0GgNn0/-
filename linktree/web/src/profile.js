import { ref } from 'vue'
import { api } from './api'

export const profile = ref({ display_name: '啊芃', bio: '', theme: 'auto' })

export async function loadProfile() {
  try {
    const data = await api.profile()
    profile.value = { ...profile.value, ...data }
    document.title = profile.value.display_name || '啊芃'
  } catch (_) { /* keep defaults */ }
}
