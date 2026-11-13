<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { listDirectory } from './api'

const loading = ref(false)
const error = ref('')
const currentPath = ref('')
const entries = ref([])
const selectedIndex = ref(0)
const panelRef = ref(null)

const selectedEntry = computed(() => entries.value[selectedIndex.value] ?? null)

async function load(path = '') {
  loading.value = true
  error.value = ''

  try {
    const response = await listDirectory(path)
    currentPath.value = response.path
    entries.value = response.entries ?? []
    selectedIndex.value = entries.value.length > 0 ? 0 : -1
    await nextTick()
    panelRef.value?.focus()
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err)
  } finally {
    loading.value = false
  }
}

function select(index) {
  if (index < 0 || index >= entries.value.length) return
  selectedIndex.value = index
}

function move(delta) {
  if (entries.value.length === 0) return
  const next = Math.min(entries.value.length - 1, Math.max(0, selectedIndex.value + delta))
  selectedIndex.value = next
  scrollSelectedIntoView()
}

function scrollSelectedIntoView() {
  nextTick(() => {
    const el = document.querySelector(`[data-row-index="${selectedIndex.value}"]`)
    el?.scrollIntoView({ block: 'nearest' })
  })
}

function openEntry(entry = selectedEntry.value) {
  if (!entry || !entry.isDir) return
  load(entry.path)
}

function onKeydown(event) {
  if (event.key === 'ArrowDown') {
    event.preventDefault()
    move(1)
  } else if (event.key === 'ArrowUp') {
    event.preventDefault()
    move(-1)
  } else if (event.key === 'Enter') {
    event.preventDefault()
    openEntry()
  }
}

function formatSize(entry) {
  if (entry.isParent) return ''
  if (entry.isDir) return '<DIR>'
  return String(entry.size)
}

onMounted(() => {
  load()
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <section class="file-panel-page p-3">
    <div class="commander-panel mx-auto" ref="panelRef" tabindex="0">
      <header class="panel-header d-flex align-items-center justify-content-between gap-3">
        <div>
          <div class="panel-title">File panel</div>
          <div class="panel-path">{{ currentPath || '...' }}</div>
        </div>
        <button class="btn btn-sm btn-outline-light" type="button" @click="load(currentPath)" :disabled="loading">
          Refresh
        </button>
      </header>

      <div v-if="error" class="alert alert-danger m-3 py-2">{{ error }}</div>
      <div v-if="loading" class="px-3 py-2 text-info">Loading...</div>

      <div class="file-table" role="listbox" aria-label="Directory entries">
        <div class="file-row file-row-head">
          <div>Name</div>
          <div class="text-end">Size</div>
          <div>Modified</div>
        </div>

        <button
          v-for="(entry, index) in entries"
          :key="entry.path + ':' + entry.name"
          class="file-row file-row-button"
          :class="{ selected: index === selectedIndex }"
          :data-row-index="index"
          type="button"
          role="option"
          :aria-selected="index === selectedIndex"
          @click="select(index)"
          @dblclick="openEntry(entry)"
        >
          <a href="#" class="entry-name" @click.prevent="entry.isDir ? openEntry(entry) : select(index)">
            <span v-if="entry.isDir" class="entry-mark">▸</span>
            <span v-else class="entry-mark">·</span>
            {{ entry.name }}
          </a>
          <span class="text-end">{{ formatSize(entry) }}</span>
          <span>{{ entry.mtime }}</span>
        </button>
      </div>

      <footer class="panel-footer">
        ↑/↓ — select, Enter/double click — open directory
      </footer>
    </div>
  </section>
</template>

<style scoped>
.file-panel-page {
  min-height: 100vh;
  background: radial-gradient(circle at top, #16213a 0, #0b1220 55%);
  color: #d7e6ff;
}

.commander-panel {
  height: calc(100vh - 2rem);
  max-width: 980px;
  border: 1px solid #315b9a;
  border-radius: 10px;
  background: #061022;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.35);
  overflow: hidden;
  outline: none;
}

.panel-header,
.panel-footer {
  background: #0d2b55;
  border-bottom: 1px solid #315b9a;
  padding: 0.75rem 1rem;
}

.panel-footer {
  border-top: 1px solid #315b9a;
  border-bottom: 0;
  color: #8fb4e8;
  font-size: 0.9rem;
}

.panel-title {
  font-weight: 700;
  letter-spacing: 0.03em;
  text-transform: uppercase;
}

.panel-path {
  color: #8fb4e8;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', monospace;
  font-size: 0.9rem;
  word-break: break-all;
}

.file-table {
  height: calc(100% - 122px);
  overflow: auto;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', monospace;
}

.file-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 120px 170px;
  gap: 1rem;
  width: 100%;
  min-height: 32px;
  align-items: center;
  padding: 0 0.9rem;
  border: 0;
  border-bottom: 1px solid rgba(49, 91, 154, 0.28);
  background: transparent;
  color: inherit;
  text-align: left;
}

.file-row-head {
  position: sticky;
  top: 0;
  z-index: 1;
  min-height: 34px;
  background: #0a1d3d;
  color: #9ec4ff;
  font-weight: 700;
}

.file-row-button:hover,
.file-row-button.selected {
  background: #134f9a;
  color: #fff;
}

.entry-name {
  overflow: hidden;
  color: inherit;
  text-decoration: none;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.entry-name:hover {
  text-decoration: underline;
}

.entry-mark {
  display: inline-block;
  width: 1.4rem;
  color: #77d0ff;
}
</style>
