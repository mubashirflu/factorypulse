<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import apiClient from '@/api/client'

interface Machine {
  id: number
  name: string
  type: string
  status: string
  location: string
}

interface Alert {
  id: number
  machine_id: number
  severity: string
  message: string
  value: number
  threshold: number
  created_at: string
}

const authStore = useAuthStore()
const router = useRouter()
const machines = ref<Machine[]>([])
const alerts = ref<Alert[]>([])
const loading = ref(true)

const navItems = [
  { label: 'Dashboard', icon: '◆', path: '/dashboard' },
  { label: 'Machines', icon: '▣', path: '/machines' },
  // { label: 'Sensors', icon: '◈', path: '/machines' },
  { label: 'Maintenance', icon: '✚', path: '/maintenance' },
  { label: 'Alerts', icon: '▲', path: '/alerts' },
    { label: 'Analytics', icon: '▤', path: '/analytics' },
]
async function fetchData() {
  try {
    const [machinesRes, alertsRes] = await Promise.all([
      apiClient.get('/machines'),
      apiClient.get('/alerts'),
    ])
    machines.value = machinesRes.data ?? []
    alerts.value = alertsRes.data ?? []
  } catch (err) {
    console.error('Failed to load dashboard data', err)
  } finally {
    loading.value = false
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

// function statusColor(status: string) {
//   if (status === 'RUNNING') return 'text-live bg-live/10 border-live/30'
//   if (status === 'WARNING') return 'text-accent bg-accent/10 border-accent/30'
//   if (status === 'CRITICAL') return 'text-critical bg-critical/10 border-critical/30'
//   return 'text-text-muted bg-surface-raised border-border'
// }
function statusColor(status: string) {
  if (status === 'RUNNING') return 'text-live bg-live/10 border-live/30'
  if (status === 'WARNING') return 'text-accent bg-accent/10 border-accent/30'
  if (status === 'CRITICAL') return 'text-critical bg-critical/10 border-critical/30'
  if (status === 'OFFLINE') return 'text-text-muted bg-surface-raised border-text-muted/30'
  return 'text-text-muted bg-surface-raised border-border'
}

function machineName(machineId: number) {
  return machines.value.find(m => m.id === machineId)?.name ?? `Machine #${machineId}`
}

onMounted(() => {
  fetchData()
  setInterval(fetchData, 5000)
})
</script>

<template>
  <div class="min-h-screen bg-base flex">
    <!-- SIDEBAR -->
    <aside class="w-60 border-r border-border flex flex-col shrink-0">
      <div class="px-6 py-6 border-b border-border">
        <h1 class="font-display text-xl font-bold text-text-primary">
          Factory<span class="text-accent">Pulse</span>
        </h1>
      </div>

      <nav class="flex-1 px-3 py-4 space-y-1">
        <!-- <a
          v-for="item in navItems"
          :key="item.label"
          href="#"
          class="flex items-center gap-3 px-3 py-2.5 rounded-md text-sm text-text-muted hover:bg-surface hover:text-text-primary transition-colors"
        >
          <span class="text-accent text-xs">{{ item.icon }}</span>
          {{ item.label }}
        </a> -->
        <router-link
  v-for="item in navItems"
  :key="item.label"
  :to="item.path"
  class="flex items-center gap-3 px-3 py-2.5 rounded-md text-sm text-text-muted hover:bg-surface hover:text-text-primary transition-colors"
>
  <span class="text-accent text-xs">{{ item.icon }}</span>
  {{ item.label }}
</router-link>
      </nav>

      <div class="p-3 border-t border-border">
        <button
          @click="handleLogout"
          class="w-full text-left px-3 py-2.5 rounded-md text-sm text-text-muted hover:bg-surface hover:text-critical transition-colors"
        >
          Log out
        </button>
      </div>
    </aside>

    <!-- MAIN CONTENT -->
    <main class="flex-1 p-8 max-w-7xl">
      <!-- CRITICAL ALERT BANNER -->
      <div
        v-if="alerts.some(a => a.severity === 'CRITICAL')"
        class="flex items-center justify-between bg-critical/10 border border-critical/30 rounded-lg px-5 py-3 mb-6"
      >
        <div class="flex items-center gap-3">
          <span class="text-critical text-lg">▲</span>
          <span class="text-critical font-medium text-sm">
            {{ alerts.filter(a => a.severity === 'CRITICAL').length }} Critical Alert{{ alerts.filter(a => a.severity === 'CRITICAL').length > 1 ? 's' : '' }}
          </span>
          <span class="text-text-muted text-sm">
            {{ machineName(alerts.find(a => a.severity === 'CRITICAL')!.machine_id) }} →
            {{ alerts.find(a => a.severity === 'CRITICAL')!.message }}
          </span>
        </div>
        <router-link
          :to="`/machines/${alerts.find(a => a.severity === 'CRITICAL')!.machine_id}`"
          class="text-critical text-sm hover:underline"
        >
          View details →
        </router-link>
      </div>

      <div class="flex items-center justify-between mb-8">
        <div>
          <h2 class="font-display text-2xl font-semibold text-text-primary">Factory overview</h2>
          <p class="text-text-muted text-sm mt-1">Live status across all connected machines</p>
        </div>
        <div class="flex items-center gap-2 font-mono text-xs text-text-muted">
          <span class="w-2 h-2 rounded-full bg-live"></span>
          Live
        </div>
      </div>

      <!-- STAT CARDS -->
      <div class="grid grid-cols-4 gap-4 mb-8">
        <div class="bg-surface border border-border rounded-lg p-5">
          <p class="text-text-muted text-xs uppercase tracking-wide font-mono">Total machines</p>
          <p class="font-display text-3xl font-bold text-text-primary mt-2">{{ machines.length }}</p>
        </div>
        <div class="bg-surface border border-border rounded-lg p-5">
          <p class="text-text-muted text-xs uppercase tracking-wide font-mono">Running</p>
          <p class="font-display text-3xl font-bold text-live mt-2">
            {{ machines.filter(m => m.status === 'RUNNING').length }}
          </p>
        </div>
        <div class="bg-surface border border-border rounded-lg p-5">
          <p class="text-text-muted text-xs uppercase tracking-wide font-mono">Warning</p>
          <p class="font-display text-3xl font-bold text-accent mt-2">
            {{ machines.filter(m => m.status === 'WARNING').length }}
          </p>
        </div>
        <div class="bg-surface border border-border rounded-lg p-5">
          <p class="text-text-muted text-xs uppercase tracking-wide font-mono">Critical</p>
          <p class="font-display text-3xl font-bold text-critical mt-2">
            {{ machines.filter(m => m.status === 'CRITICAL').length }}
          </p>
        </div>
      </div>

      <!-- MACHINES TABLE -->
      <div class="bg-surface border border-border rounded-lg overflow-hidden">
        <div class="px-5 py-4 border-b border-border">
          <h3 class="font-display font-semibold text-text-primary">Machines</h3>
        </div>

        <div v-if="loading" class="p-8 text-center text-text-muted text-sm">Loading machines…</div>

        <div v-else-if="machines.length === 0" class="p-8 text-center text-text-muted text-sm">
          No machines yet. Add one from the Machines page.
        </div>

        <table v-else class="w-full text-sm">
          <thead>
            <tr class="text-text-muted text-xs uppercase tracking-wide font-mono border-b border-border">
              <th class="text-left px-5 py-3 font-normal">Name</th>
              <th class="text-left px-5 py-3 font-normal">Type</th>
              <th class="text-left px-5 py-3 font-normal">Location</th>
              <th class="text-left px-5 py-3 font-normal">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="machine in machines"
              :key="machine.id"
              @click="router.push(`/machines/${machine.id}`)"
              class="border-b border-border last:border-0 hover:bg-surface-raised transition-colors cursor-pointer"
            >
              <td class="px-5 py-3 text-text-primary font-medium">{{ machine.name }}</td>
              <td class="px-5 py-3 text-text-muted">{{ machine.type }}</td>
              <td class="px-5 py-3 text-text-muted">{{ machine.location || '—' }}</td>
              <td class="px-5 py-3">
                <span :class="['inline-block px-2.5 py-1 rounded text-xs font-mono border', statusColor(machine.status)]">
                  {{ machine.status }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </main>
  </div>
</template>