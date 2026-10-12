import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api'
import { config, serverUrl } from '../config'

/** Whether the browser holds a panel session. */
export const useSession = defineStore('session', () => {
  const loggedIn = ref(config.loggedIn)

  /** Throws ApiError when the server refuses the credentials. */
  async function login(username: string, password: string, twoFactorCode: string) {
    await api.post('login', { username, password, twoFactorCode })
    loggedIn.value = true
  }

  async function logout() {
    // The server clears the session cookie and answers with a redirect we do not follow.
    await fetch(serverUrl('logout'), { credentials: 'same-origin', redirect: 'manual' }).catch(() => {})
    loggedIn.value = false
  }

  return { loggedIn, login, logout }
})
