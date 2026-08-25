<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import apiClient from '@/api/client'
import { useAuthStore } from '@/stores/auth'

interface Machine {
  id: number
  name: string
  type: string
  status: string
  location: string
}

const router = useRouter()
const authStore=useAuthStore()
const machines = ref<Machine[]>([])
const loading = ref(true)
const showForm = ref(false)

const newName = ref('')
const newType = ref('')
const newLocation = ref('')
const submitting = ref(false)
const errorMessage = ref('')

async function fetchMachines() {
  loading.value = true
  try {
    const res = await apiClient.get('/machines')
    machines.value = res.data ?? []
  } catch (err) {
    console.error('Failed to load machines', err)
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  errorMessage.value = ''
  submitting.value = true
  try {
    await apiClient.post('/machines', {
      name: newName.value,
      type: newType.value,
      location: newLocation.value,
    })
    newName.value = ''
    newType.value = ''
    newLocation.value = ''
    showForm.value = false
    await fetchMachines()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Could not create machine'
  } finally {
    submitting.value = false
  }
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

// onMounted(fetchMachines)
onMounted(() => {
  fetchMachines()
  setInterval(fetchMachines, 5000)
})
</script>

<template>
  <div class="min-h-screen bg-base p-8">
    <button @click="router.push('/dashboard')" class="text-text-muted text-sm hover:text-text-primary mb-6">
      ← Back to dashboard
    </button>

    <div class="flex items-center justify-between mb-8">
      <div>
        <h2 class="font-display text-2xl font-semibold text-text-primary">Machines</h2>
        <p class="text-text-muted text-sm mt-1">Manage your factory equipment</p>
      </div>
      <!-- <button
        @click="showForm = !showForm"
        class="bg-accent hover:bg-accent-dim text-base font-semibold px-4 py-2.5 rounded-md text-sm transition-colors"
      >
        {{ showForm ? 'Cancel' : '+ Add machine' }}
      </button> -->
        <button
   v-if="authStore.isAdmin()"
  @click="showForm = !showForm"
    class="bg-accent hover:bg-accent-dim text-base font-semibold px-4 py-2.5 rounded-md text-sm transition-colors"
  >
    {{ showForm ? 'Cancel' : '+ Add machine' }}
  </button>
    </div>

    <div v-if="showForm" class="bg-surface border border-border rounded-lg p-6 mb-6">
      <h3 class="font-display font-semibold text-text-primary mb-4">New machine</h3>
      <form @submit.prevent="handleCreate" class="grid grid-cols-3 gap-4">
        <div>
          <label class="block text-text-muted text-xs uppercase tracking-wide mb-2 font-mono">Name</label>
          <input
            v-model="newName"
            required
            placeholder="CNC-02"
            class="w-full px-3 py-2.5 rounded-md bg-base border border-border text-text-primary placeholder-text-muted outline-none focus:border-accent"
          />
        </div>
        <div>
          <label class="block text-text-muted text-xs uppercase tracking-wide mb-2 font-mono">Type</label>
          <input
            v-model="newType"
            required
            placeholder="CNC"
            class="w-full px-3 py-2.5 rounded-md bg-base border border-border text-text-primary placeholder-text-muted outline-none focus:border-accent"
          />
        </div>
        <div>
          <label class="block text-text-muted text-xs uppercase tracking-wide mb-2 font-mono">Location</label>
          <input
            v-model="newLocation"
            placeholder="Floor A"
            class="w-full px-3 py-2.5 rounded-md bg-base border border-border text-text-primary placeholder-text-muted outline-none focus:border-accent"
          />
        </div>
        <div class="col-span-3 flex items-center gap-4">
          <button
            type="submit"
            :disabled="submitting"
            class="bg-accent hover:bg-accent-dim disabled:opacity-50 text-base font-semibold px-4 py-2.5 rounded-md text-sm transition-colors"
          >
            {{ submitting ? 'Adding…' : 'Add machine' }}
          </button>
          <p v-if="errorMessage" class="text-critical text-sm">{{ errorMessage }}</p>
        </div>
      </form>
    </div>

    <div v-if="loading" class="text-text-muted text-sm">Loading…</div>
    <div v-else-if="machines.length === 0" class="text-text-muted text-sm">No machines yet.</div>
    <div v-else class="grid grid-cols-3 gap-4">
      <div
        v-for="machine in machines"
        :key="machine.id"
        @click="router.push(`/machines/${machine.id}`)"
        class="bg-surface border border-border rounded-lg p-5 cursor-pointer hover:border-accent/50 transition-colors"
      >
        <div class="flex items-center justify-between mb-3">
          <h4 class="font-display font-semibold text-text-primary">{{ machine.name }}</h4>
          <span :class="['inline-block px-2 py-0.5 rounded text-xs font-mono border', statusColor(machine.status)]">
            {{ machine.status }}
          </span>
        </div>
        <p class="text-text-muted text-sm">{{ machine.type }}</p>
        <p class="text-text-muted text-xs mt-1">{{ machine.location || 'No location set' }}</p>
      </div>
    </div>
  </div>
</template>