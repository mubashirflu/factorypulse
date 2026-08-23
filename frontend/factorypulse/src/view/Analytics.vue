<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import apiClient from '@/api/client'
import { Chart, registerables } from 'chart.js'

Chart.register(...registerables)

interface MachineAnalytics {
  machine_id: number
  machine_name: string
  total_alerts: number
  critical_count: number
  avg_repair_minutes: number
}

const router = useRouter()
const data = ref<MachineAnalytics[]>([])
const loading = ref(true)
const canvasRef = ref<HTMLCanvasElement | null>(null)
let chartInstance: Chart | null = null

async function fetchAnalytics() {
  try {
    const res = await apiClient.get('/analytics')
    data.value = res.data ?? []
    renderChart()
  } catch (err) {
    console.error('Failed to load analytics', err)
  } finally {
    loading.value = false
  }
}

function renderChart() {
  if (!canvasRef.value) return
  if (chartInstance) chartInstance.destroy()

  chartInstance = new Chart(canvasRef.value, {
    type: 'bar',
    data: {
      labels: data.value.map(d => d.machine_name),
      datasets: [
        {
          label: 'Total alerts',
          data: data.value.map(d => d.total_alerts),
          backgroundColor: '#e8a33d',
          borderRadius: 4,
        },
        {
          label: 'Critical alerts',
          data: data.value.map(d => d.critical_count),
          backgroundColor: '#e4572e',
          borderRadius: 4,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { labels: { color: '#8a8f98' } },
      },
      scales: {
        x: { ticks: { color: '#8a8f98' }, grid: { display: false } },
        y: { ticks: { color: '#8a8f98' }, grid: { color: '#262b31' }, beginAtZero: true },
      },
    },
  })
}

function formatMinutes(mins: number) {
  if (mins === 0) return 'No data'
  if (mins < 60) return `${mins.toFixed(0)} min`
  return `${(mins / 60).toFixed(1)} hrs`
}

onMounted(fetchAnalytics)
</script>

<template>
  <div class="min-h-screen bg-base p-8">
    <button @click="router.push('/dashboard')" class="text-text-muted text-sm hover:text-text-primary mb-6">
      ← Back to dashboard
    </button>

    <div class="mb-8">
      <h2 class="font-display text-2xl font-semibold text-text-primary">Analytics</h2>
      <p class="text-text-muted text-sm mt-1">Alert frequency and repair performance by machine</p>
    </div>

    <div v-if="loading" class="text-text-muted text-sm">Loading…</div>

    <template v-else>
      <!-- CHART -->
      <div class="bg-surface border border-border rounded-lg p-5 mb-6">
        <h3 class="font-display font-semibold text-text-primary mb-4">Alerts by machine</h3>
        <div class="relative h-72">
          <canvas ref="canvasRef"></canvas>
        </div>
      </div>

      <!-- MTTR TABLE -->
      <div class="bg-surface border border-border rounded-lg overflow-hidden">
        <div class="px-5 py-4 border-b border-border">
          <h3 class="font-display font-semibold text-text-primary">Mean time to repair (MTTR)</h3>
        </div>
        <table class="w-full text-sm">
          <thead>
            <tr class="text-text-muted text-xs uppercase tracking-wide font-mono border-b border-border">
              <th class="text-left px-5 py-3 font-normal">Machine</th>
              <th class="text-left px-5 py-3 font-normal">Total alerts</th>
              <th class="text-left px-5 py-3 font-normal">Critical</th>
              <th class="text-left px-5 py-3 font-normal">Avg repair time</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in data"
              :key="row.machine_id"
              class="border-b border-border last:border-0 hover:bg-surface-raised transition-colors"
            >
              <td class="px-5 py-3 text-text-primary font-medium">{{ row.machine_name }}</td>
              <td class="px-5 py-3 text-text-muted font-mono">{{ row.total_alerts }}</td>
              <td class="px-5 py-3 font-mono" :class="row.critical_count > 0 ? 'text-critical' : 'text-text-muted'">
                {{ row.critical_count }}
              </td>
              <td class="px-5 py-3 text-text-muted font-mono">{{ formatMinutes(row.avg_repair_minutes) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>