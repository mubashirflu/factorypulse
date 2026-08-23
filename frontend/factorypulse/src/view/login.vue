<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const email = ref('')
const password = ref('')
const errorMessage = ref('')
const loading = ref(false)

const authStore = useAuthStore()
const router = useRouter()

async function handleLogin() {
  errorMessage.value = ''
  loading.value = true
  try {
    await authStore.login(email.value, password.value)
    router.push('/dashboard')
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Login failed. Check your credentials.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-base flex">
    <!-- LEFT PANEL: Branding + signature pulse element -->
    <div class="hidden lg:flex lg:w-1/2 flex-col justify-between p-16 border-r border-border relative overflow-hidden">
      <div>
        <div class="flex items-center gap-3">
          <div class="w-2.5 h-2.5 rounded-full bg-live"></div>
          <span class="text-text-muted text-xs tracking-[0.2em] uppercase font-mono">System online</span>
        </div>
        <h1 class="font-display text-5xl font-bold text-text-primary mt-8 leading-tight">
          Factory<span class="text-accent">Pulse</span>
        </h1>
        <p class="text-text-muted mt-4 max-w-sm leading-relaxed">
          Real-time machine telemetry, maintenance workflows and downtime
          analysis for the industrial floor.
        </p>
      </div>

      <!-- Signature: animated pulse line -->
      <div class="relative">
        <svg viewBox="0 0 400 100" class="w-full h-24">
          <polyline
            points="0,50 60,50 80,50 95,15 110,85 125,50 160,50 180,50 195,25 210,75 225,50 400,50"
            fill="none"
            stroke="#4fd1c5"
            stroke-width="2"
            stroke-linejoin="round"
            stroke-linecap="round"
            class="pulse-line"
          />
        </svg>
        <div class="flex justify-between mt-4 font-mono text-xs text-text-muted">
          <span>TEMP <span class="text-live">72.4°C</span></span>
          <span>VIB <span class="text-live">3.1mm/s</span></span>
          <span>STATUS <span class="text-live">RUNNING</span></span>
        </div>
      </div>
    </div>

    <!-- RIGHT PANEL: Login form -->
    <div class="w-full lg:w-1/2 flex items-center justify-center p-8">
      <div class="w-full max-w-sm">
        <div class="lg:hidden mb-10">
          <h1 class="font-display text-3xl font-bold text-text-primary">
            Factory<span class="text-accent">Pulse</span>
          </h1>
        </div>

        <h2 class="font-display text-2xl font-semibold text-text-primary mb-1">Sign in</h2>
        <p class="text-text-muted text-sm mb-8">Access your factory dashboard</p>

        <form @submit.prevent="handleLogin" class="space-y-5">
          <div>
            <label class="block text-text-muted text-xs tracking-wide uppercase mb-2 font-mono">Email</label>
            <input
              v-model="email"
              type="email"
              required
              placeholder="name@company.com"
              class="w-full px-4 py-3 rounded-md bg-surface border border-border text-text-primary placeholder-text-muted outline-none focus:border-accent transition-colors"
            />
          </div>

          <div>
            <label class="block text-text-muted text-xs tracking-wide uppercase mb-2 font-mono">Password</label>
            <input
              v-model="password"
              type="password"
              required
              placeholder="••••••••"
              class="w-full px-4 py-3 rounded-md bg-surface border border-border text-text-primary placeholder-text-muted outline-none focus:border-accent transition-colors"
            />
          </div>

          <p v-if="errorMessage" class="text-critical text-sm">{{ errorMessage }}</p>

          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-accent hover:bg-accent-dim disabled:opacity-50 text-base font-semibold py-3 rounded-md transition-colors mt-2"
          >
            {{ loading ? 'Signing in…' : 'Sign in' }}
          </button>
          <p class="text-text-muted text-sm mt-6 text-center">
  Don't have an account?
  <router-link to="/register" class="text-accent hover:underline">Create one</router-link>
</p>
          
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pulse-line {
  stroke-dasharray: 500;
  stroke-dashoffset: 500;
  animation: draw 3s ease-in-out infinite;
}
@keyframes draw {
  0% { stroke-dashoffset: 500; }
  50% { stroke-dashoffset: 0; }
  100% { stroke-dashoffset: -500; }
}
</style>