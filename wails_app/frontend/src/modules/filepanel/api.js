import { backendBaseUrl } from '../../appConfig'

async function getJson(path) {
  const response = await fetch(`${backendBaseUrl}${path}`)
  if (!response.ok) {
    const text = await response.text()
    throw new Error(text || `HTTP ${response.status}`)
  }
  return response.json()
}

export function getFilePanelConfig() {
  return getJson('/api/filepanel/config')
}

export function listDirectory(path) {
  const query = path ? `?path=${encodeURIComponent(path)}` : ''
  return getJson(`/api/filepanel/list${query}`)
}
