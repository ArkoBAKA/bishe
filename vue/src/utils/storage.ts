const TOKEN_KEY = 'token'
const USER_KEY = 'user'

export const getToken = () => window.localStorage.getItem(TOKEN_KEY) || ''

export const setToken = (token: string) => {
  if (!token) {
    window.localStorage.removeItem(TOKEN_KEY)
    return
  }
  window.localStorage.setItem(TOKEN_KEY, token)
}

export const getUser = <T = unknown>(): T | null => {
  const raw = window.localStorage.getItem(USER_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as T
  } catch (e) {
    return null
  }
}

export const setUser = (user: unknown) => {
  if (!user) {
    window.localStorage.removeItem(USER_KEY)
    return
  }
  window.localStorage.setItem(USER_KEY, JSON.stringify(user))
}

export const clearAuth = () => {
  window.localStorage.removeItem(TOKEN_KEY)
  window.localStorage.removeItem(USER_KEY)
}

