import { defineStore } from 'pinia'
import { clearAuth, getToken, getUser, setToken, setUser } from '@/utils/storage'
import type { UserSummary } from '@/types/api'
import { usersApi } from '@/apis'

export const useAuthStore = defineStore('auth', {
  state: (): { token: string; user: UserSummary | null } => ({
    token: getToken(),
    user: getUser<UserSummary>()
  }),
  getters: {
    isAuthed: (state): boolean => !!state.token,
    isAdmin: (state): boolean => state.user?.role === 'admin'
  },
  actions: {
    async login(account: string, password: string) {
      const data = await usersApi.login({ account, password })
      this.token = data.token
      this.user = data.user
      setToken(data.token)
      setUser(data.user)
      return data
    },
    async fetchMe() {
      const data = await usersApi.getMe()
      this.user = data
      setUser(data)
      return data
    },
    logout() {
      this.token = ''
      this.user = null
      clearAuth()
    }
  }
})
