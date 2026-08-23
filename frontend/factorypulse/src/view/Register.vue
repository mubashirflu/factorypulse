<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const name = ref('')
const email = ref('')
const password = ref('')
const errorMessage = ref('')
const loading = ref(false)

const authStore = useAuthStore()
const router = useRouter()

async function handleRegister() {
  errorMessage.value = ''
  loading.value = true
  try {
    await authStore.register(name.value, email.value, password.value)
    // Register ke baad seedha login kar do, taaki user ko dobara form na bharna pade
    await authStore.login(email.value, password.value)
    router.push('/dashboard')
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Registration failed. Try again.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-base flex">
    <div class="hidden lg:flex lg:w-1/2 flex-col justify-between p-16 border-r border-border">
      <div>
        <div class="flex items-center gap-3">
          <div class="w-2.5 h-2.5 rounded-full bg-live"></div>
          <span class="text-text-muted text-xs tracking-[0.2em] uppercase font-mono">System online</span>
        </div>
        <h1 class="font-display text-5xl font-bold text-text-primary mt-8 leading-tight">
          Factory<span class="text-accent">Pulse</span>
        </h1>
        <p class="text-text-muted mt-4 max-w-sm leading-relaxed">
          Create an account to monitor machines, manage maintenance jobs
          and track factory-wide analytics.
        </p>
      </div>
    </div>

    <div class="w-full lg:w-1/2 flex items-center justify-center p-8">
      <div class="w-full max-w-sm">
        <div class="lg:hidden mb-10">
          <h1 class="font-display text-3xl font-bold text-text-primary">
            Factory<span class="text-accent">Pulse</span>
          </h1>
        </div>

        <h2 class="font-display text-2xl font-semibold text-text-primary mb-1">Create account</h2>
        <p class="text-text-muted text-sm mb-8">Get access to the factory dashboard</p>

        <form @submit.prevent="handleRegister" class="space-y-5">
          <div>
            <label class="block text-text-muted text-xs tracking-wide uppercase mb-2 font-mono">Name</label>
            <input
              v-model="name"
              type="text"
              required
              placeholder="Your full name"
              class="w-full px-4 py-3 rounded-md bg-surface border border-border text-text-primary placeholder-text-muted outline-none focus:border-accent transition-colors"
            />
          </div>

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
              placeholder="At least 8 characters"
              class="w-full px-4 py-3 rounded-md bg-surface border border-border text-text-primary placeholder-text-muted outline-none focus:border-accent transition-colors"
            />
          </div>

          <p v-if="errorMessage" class="text-critical text-sm">{{ errorMessage }}</p>

          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-accent hover:bg-accent-dim disabled:opacity-50 text-base font-semibold py-3 rounded-md transition-colors mt-2"
          >
            {{ loading ? 'Creating account…' : 'Create account' }}
          </button>
        </form>

        <p class="text-text-muted text-sm mt-6 text-center">
          Already have an account?
          <router-link to="/login" class="text-accent hover:underline">Sign in</router-link>
        </p>
      </div>
    </div>
  </div>
</template>