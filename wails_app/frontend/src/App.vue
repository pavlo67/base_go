<script setup>
import FilePanel from './modules/filepanel/FilePanel.vue'
import { activeModule } from './appConfig'
</script>

<template>
  <main class="app-shell">
    <FilePanel v-if="activeModule === 'filepanel'" />
    <section v-else class="p-4">
      Unknown module: <strong>{{ activeModule }}</strong>
    </section>
  </main>
</template>

<style>
html,
body,
#app {
  height: 100%;
  margin: 0;
}

body {
  background: #0b1220;
}

.app-shell {
  min-height: 100%;
}
</style>
