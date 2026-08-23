<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import apiClient from '@/api/client'

interface Alert {
  id: number
  machine_id: number
  severity: string
  message: string
  value: number
  threshold: number
  resolved: boolean
  created_at: string
}

interface Machine {
  id: number
  name: string
}

const router = useRouter()
const allAlerts = ref<Alert[]>([])
const machines = ref<Machine[]>([])
const loading = ref(true)
const filter = ref<'all' | 'active' | 'resolved'>('active')

async function fetchData() {
  loading.value = true
  try {
    const [alertsRes, machinesRes] = await Promise.all([
      apiClient.get('/alerts/all'),
      apiClient.get('/machines'),
    ])
    allAlerts.value = alertsRes.data ?? []
    machines.value = machinesRes.data ?? []
  } catch (err) {
    console.error('Failed to load alerts', err)
  } finally {
    loading.value = false
  }
}

async function resolveAlert(id: number) {
  try {
    await apiClient.patch(`/alerts/${id}/resolve`)
    await fetchData()
  } catch (err) {
    console.error('Failed to resolve alert', err)
  }
}

function machineName(machineId: number) {
  return machines.value.find(m => m.id === machineId)?.name ?? `Machine #${machineId}`
}

function filteredAlerts() {
  if (filter.value === 'active') return allAlerts.value.filter(a => !a.resolved)
  if (filter.value === 'resolved') return allAlerts.value.filter(a => a.resolved)
  return allAlerts.value
}

function severityStyle(severity: string) {
  if (severity === 'CRITICAL') return 'text-critical bg-critical/10 border-critical/30'
  return 'text-accent bg-accent/10 border-accent/30'
}

function timeAgo(dateStr: string) {
  const diffMs = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diffMs / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  return `${Math.floor(hrs / 24)}d ago`
}

onMounted(fetchData)
</script>

<template>
  <div class="min-h-screen bg-base p-8">
    <button @click="router.push('/dashboard')" class="text-text-muted text-sm hover:text-text-primary mb-6">
      ← Back to dashboard
    </button>

    <div class="flex items-center justify-between mb-8">
      <div>
        <h2 class="font-display text-2xl font-semibold text-text-primary">Alerts</h2>
        <p class="text-text-muted text-sm mt-1">Machine warnings and critical events</p>
      </div>

      <div class="flex gap-2">
        <button
          v-for="f in ['active', 'all', 'resolved']"
          :key="f"
          @click="filter = f as any"
          :class="[
            'px-3 py-1.5 rounded-md text-xs font-mono capitalize transition-colors',
            filter === f ? 'bg-accent text-base font-semibold' : 'bg-surface text-text-muted border border-border'
          ]"
        >
          {{ f }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-text-muted text-sm">Loading…</div>
    <div v-else-if="filteredAlerts().length === 0" class="text-text-muted text-sm">No alerts here.</div>

    <div v-else class="space-y-3">
      <div
        v-for="alert in filteredAlerts()"
        :key="alert.id"
        class="bg-surface border border-border rounded-lg p-4 flex items-center justify-between"
      >
        <div class="flex items-center gap-4">
          <span :class="['px-2 py-1 rounded text-xs font-mono border', severityStyle(alert.severity)]">
            {{ alert.severity }}
          </span>
          <div>
            <p class="text-text-primary text-sm font-medium">
              {{ machineName(alert.machine_id) }}
              <span class="text-text-muted font-normal">→ {{ alert.message }}</span>
            </p>
            <p class="text-text-muted text-xs mt-1 font-mono">{{ timeAgo(alert.created_at) }}</p>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <router-link
            :to="`/machines/${alert.machine_id}`"
            class="text-accent text-xs hover:underline"
          >
            View machine
          </router-link>
          <button
            v-if="!alert.resolved"
            @click="resolveAlert(alert.id)"
            class="text-live text-xs hover:underline"
          >
            Mark resolved
          </button>
          <span v-else class="text-text-muted text-xs">✓ Resolved</span>
        </div>
      </div>
    </div>
  </div>
</template>