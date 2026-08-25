<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import apiClient from '@/api/client'

interface Job {
  id: number
  machine_id: number
  title: string
  description: string
  status: string
  assigned_to: number | null
  created_at: string
}

interface Machine {
  id: number
  name: string
}

const router = useRouter()
const jobs = ref<Job[]>([])
const machines = ref<Machine[]>([])
const loading = ref(true)
const showForm = ref(false)

const newTitle = ref('')
const newDescription = ref('')
const newMachineId = ref<number | null>(null)
const submitting = ref(false)
const errorMessage = ref('')

const columns = [
  { status: 'OPEN', label: 'Open', color: 'border-t-text-muted' },
  { status: 'ASSIGNED', label: 'Assigned', color: 'border-t-accent' },
  { status: 'IN_PROGRESS', label: 'In progress', color: 'border-t-live' },
  { status: 'COMPLETED', label: 'Completed', color: 'border-t-live' },
]

async function fetchData() {
  loading.value = true
  try {
    const [jobsRes, machinesRes] = await Promise.all([
      apiClient.get('/maintenance'),
      apiClient.get('/machines'),
    ])
    jobs.value = jobsRes.data ?? []
    machines.value = machinesRes.data ?? []
  } catch (err) {
    console.error('Failed to load maintenance data', err)
  } finally {
    loading.value = false
  }
}

function jobsForStatus(status: string) {
  return jobs.value.filter(j => j.status === status)
}

function machineName(machineId: number) {
  return machines.value.find(m => m.id === machineId)?.name ?? `Machine #${machineId}`
}

async function handleCreate() {
  errorMessage.value = ''
  if (!newMachineId.value) {
    errorMessage.value = 'Select a machine'
    return
  }
  submitting.value = true
  try {
    await apiClient.post('/maintenance', {
      machine_id: newMachineId.value,
      title: newTitle.value,
      description: newDescription.value,
    })
    newTitle.value = ''
    newDescription.value = ''
    newMachineId.value = null
    showForm.value = false
    await fetchData()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Could not create job'
  } finally {
    submitting.value = false
  }
}

async function advanceStatus(job: Job) {
  const flow: Record<string, string> = {
    OPEN: 'ASSIGNED',
    ASSIGNED: 'IN_PROGRESS',
    IN_PROGRESS: 'COMPLETED',
  }
  const nextStatus = flow[job.status]
  if (!nextStatus) return

  try {
    if (job.status === 'OPEN') {
      await apiClient.patch(`/maintenance/${job.id}/assign`)
    } else {
      await apiClient.patch(`/maintenance/${job.id}/status`, { status: nextStatus })
    }
    await fetchData()
  } catch (err) {
    console.error('Failed to update job', err)
  }
}

// onMounted(fetchData)
onMounted(() => {
  fetchData()
  setInterval(fetchData, 5000)
})
</script>

<template>
  <div class="min-h-screen bg-base p-8">
    <button @click="router.push('/dashboard')" class="text-text-muted text-sm hover:text-text-primary mb-6">
      ← Back to dashboard
    </button>

    <div class="flex items-center justify-between mb-8">
      <div>
        <h2 class="font-display text-2xl font-semibold text-text-primary">Maintenance</h2>
        <p class="text-text-muted text-sm mt-1">Track repair jobs across the factory</p>
      </div>
      <button
        @click="showForm = !showForm"
        class="bg-accent hover:bg-accent-dim text-base font-semibold px-4 py-2.5 rounded-md text-sm transition-colors"
      >
        {{ showForm ? 'Cancel' : '+ New job' }}
      </button>
    </div>

    <!-- NEW JOB FORM -->
    <div v-if="showForm" class="bg-surface border border-border rounded-lg p-6 mb-6">
      <h3 class="font-display font-semibold text-text-primary mb-4">New maintenance job</h3>
      <form @submit.prevent="handleCreate" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-text-muted text-xs uppercase tracking-wide mb-2 font-mono">Machine</label>
            <select
              v-model="newMachineId"
              class="w-full px-3 py-2.5 rounded-md bg-base border border-border text-text-primary outline-none focus:border-accent"
            >
              <option :value="null" disabled>Select machine</option>
              <option v-for="m in machines" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-text-muted text-xs uppercase tracking-wide mb-2 font-mono">Title</label>
            <input
              v-model="newTitle"
              required
              placeholder="Inspect spindle bearing"
              class="w-full px-3 py-2.5 rounded-md bg-base border border-border text-text-primary placeholder-text-muted outline-none focus:border-accent"
            />
          </div>
        </div>
        <div>
          <label class="block text-text-muted text-xs uppercase tracking-wide mb-2 font-mono">Description</label>
          <textarea
            v-model="newDescription"
            rows="2"
            placeholder="Optional details..."
            class="w-full px-3 py-2.5 rounded-md bg-base border border-border text-text-primary placeholder-text-muted outline-none focus:border-accent resize-none"
          ></textarea>
        </div>
        <div class="flex items-center gap-4">
          <button
            type="submit"
            :disabled="submitting"
            class="bg-accent hover:bg-accent-dim disabled:opacity-50 text-base font-semibold px-4 py-2.5 rounded-md text-sm transition-colors"
          >
            {{ submitting ? 'Creating…' : 'Create job' }}
          </button>
          <p v-if="errorMessage" class="text-critical text-sm">{{ errorMessage }}</p>
        </div>
      </form>
    </div>

    <!-- KANBAN BOARD -->
    <div v-if="loading" class="text-text-muted text-sm">Loading…</div>
    <div v-else class="grid grid-cols-4 gap-4">
      <div v-for="col in columns" :key="col.status">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-text-primary text-sm font-medium">{{ col.label }}</h3>
          <span class="text-text-muted text-xs font-mono">{{ jobsForStatus(col.status).length }}</span>
        </div>

        <div :class="['space-y-3 border-t-2 pt-3', col.color]">
          <div
            v-for="job in jobsForStatus(col.status)"
            :key="job.id"
            class="bg-surface border border-border rounded-lg p-4"
          >
            <p class="text-text-primary text-sm font-medium mb-1">{{ job.title }}</p>
            <p class="text-text-muted text-xs mb-3">{{ machineName(job.machine_id) }}</p>
            <p v-if="job.description" class="text-text-muted text-xs mb-3 line-clamp-2">{{ job.description }}</p>

            <button
              v-if="job.status !== 'COMPLETED'"
              @click="advanceStatus(job)"
              class="text-accent text-xs hover:underline"
            >
              {{ job.status === 'OPEN' ? 'Assign to me' : 'Advance status →' }}
            </button>
          </div>

          <p v-if="jobsForStatus(col.status).length === 0" class="text-text-muted text-xs italic">No jobs</p>
        </div>
      </div>
    </div>
  </div>
</template>