<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import apiClient from '@/api/client'
import { Chart, registerables } from 'chart.js'

Chart.register(...registerables)

interface Reading {
  id: number
  temperature: number
  vibration: number
  pressure: number
  recorded_at: string
}

const route = useRoute()
const router = useRouter()
const machineId = route.params.id as string

const latest = ref<Reading | null>(null)
const history = ref<Reading[]>([])
const canvasRef = ref<HTMLCanvasElement | null>(null)
let chartInstance: Chart | null = null
let pollInterval: ReturnType<typeof setInterval> | null = null

async function fetchData() {
  try {
    const [latestRes, historyRes] = await Promise.all([
      apiClient.get(`/machines/${machineId}/latest-reading`),
      apiClient.get(`/machines/${machineId}/history?limit=30`),
    ])
    latest.value = latestRes.data
    history.value = historyRes.data.reverse()
    updateChart()
  } catch (err) {
    console.error('Failed to fetch reading data', err)
  }
}

function updateChart() {
  if (!canvasRef.value) return

  const labels = history.value.map(r =>
    new Date(r.recorded_at).toLocaleTimeString('en-GB', { hour12: false })
  )
  const vibrationData = history.value.map(r => r.vibration)

  if (chartInstance) {
    chartInstance.data.labels = [...labels]
    const dataset = chartInstance.data.datasets[0]
    if (dataset) dataset.data = [...vibrationData]
    chartInstance.update()
    return
  }

  chartInstance = new Chart(canvasRef.value, {
    type: 'line',
    data: {
      labels,
      datasets: [
        {
          label: 'Vibration (mm/s)',
          data: vibrationData,
          borderColor: '#4fd1c5',
          backgroundColor: 'rgba(79, 209, 197, 0.1)',
          fill: true,
          tension: 0.3,
          pointRadius: 2,
          pointBackgroundColor: '#4fd1c5',
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      animation: { duration: 300 },
      plugins: { legend: { display: false } },
      scales: {
        x: { ticks: { color: '#8a8f98', maxTicksLimit: 6 }, grid: { color: '#262b31' } },
        y: { ticks: { color: '#8a8f98' }, grid: { color: '#262b31' }, beginAtZero: false },
      },
    },
  })
}

onMounted(() => {
  fetchData()
  pollInterval = setInterval(fetchData, 3000)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
  if (chartInstance) chartInstance.destroy()
})
</script>

<template>
  <div class="min-h-screen bg-base p-8">
    <button @click="router.push('/machines')" class="text-text-muted text-sm hover:text-text-primary mb-6">
      ← Back to machines
    </button>

    <div class="flex items-center justify-between mb-8">
      <h2 class="font-display text-2xl font-semibold text-text-primary">Machine #{{ machineId }}</h2>
      <div class="flex items-center gap-2 font-mono text-xs text-text-muted">
        <span class="w-2 h-2 rounded-full bg-live animate-pulse"></span>
        Live · updating every 3s
      </div>
    </div>

    <div class="grid grid-cols-3 gap-4 mb-8">
      <div class="bg-surface border border-border rounded-lg p-5">
        <p class="text-text-muted text-xs uppercase tracking-wide font-mono">Temperature</p>
        <p class="font-mono text-3xl font-semibold text-text-primary mt-2">
          {{ latest?.temperature ?? '—' }}<span class="text-lg text-text-muted">°C</span>
        </p>
      </div>
      <div class="bg-surface border border-border rounded-lg p-5">
        <p class="text-text-muted text-xs uppercase tracking-wide font-mono">Vibration</p>
        <p class="font-mono text-3xl font-semibold text-live mt-2">
          {{ latest?.vibration ?? '—' }}<span class="text-lg text-text-muted">mm/s</span>
        </p>
      </div>
      <div class="bg-surface border border-border rounded-lg p-5">
        <p class="text-text-muted text-xs uppercase tracking-wide font-mono">Pressure</p>
        <p class="font-mono text-3xl font-semibold text-text-primary mt-2">
          {{ latest?.pressure ?? '—' }}<span class="text-lg text-text-muted">bar</span>
        </p>
      </div>
    </div>

    <div class="bg-surface border border-border rounded-lg p-5">
      <h3 class="font-display font-semibold text-text-primary mb-4">Vibration trend</h3>
      <div class="relative h-64">
        <canvas ref="canvasRef"></canvas>
      </div>
    </div>
  </div>
</template>