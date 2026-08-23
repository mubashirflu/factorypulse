// import { defineStore } from 'pinia'
// import { ref } from 'vue'
// import apiClient from '@/api/client'

// export const useAuthStore = defineStore('auth', () => {
//   const token = ref<string | null>(localStorage.getItem('token'))
//   const isLoggedIn = ref<boolean>(!!token.value)

//   async function login(email: string, password: string) {
//     const response = await apiClient.post('/auth/login', {
//       email,
//       password,
//     })

//     token.value = response.data.token
//     localStorage.setItem('token', response.data.token)
//     isLoggedIn.value = true
//   }

//   async function register(
//     name: string,
//     email: string,
//     password: string
//   ) {
//     await apiClient.post('/auth/register', {
//       name,
//       email,
//       password,
//     })
//   }

//   function logout() {
//     token.value = null
//     localStorage.removeItem('token')
//     isLoggedIn.value = false
//   }

//   return {
//     token,
//     isLoggedIn,
//     login,
//     register,
//     logout,
//   }
// })

import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const isLoggedIn = ref<boolean>(!!token.value)
  const role = ref<string | null>(localStorage.getItem('role'))

  async function login(email: string, password: string) {
    const response = await apiClient.post('/auth/login', { email, password })
    token.value = response.data.token
    localStorage.setItem('token', response.data.token)
    isLoggedIn.value = true

    // Role fetch karo /me se
    const meResponse = await apiClient.get('/me')
    role.value = meResponse.data.role
    localStorage.setItem('role', meResponse.data.role)
  }

  async function register(name: string, email: string, password: string) {
    await apiClient.post('/auth/register', { name, email, password })
  }

  function logout() {
    token.value = null
    role.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('role')
    isLoggedIn.value = false
  }

  function isAdmin() {
    return role.value === 'admin'
  }

  return { token, isLoggedIn, role, login, register, logout, isAdmin }
})