import { defineStore } from 'pinia'
import { ref } from 'vue'
import { serverUrl, config } from '../config'

/** Whether the browser holds a panel session. */
export const useSession = defineStore('session', () => {
  const loggedIn = ref(config.loggedIn)

  async function logout() {
    // The server clears the session cookie and answers with a redirect we do not follow.
    await fetch(serverUrl('logout'), { credentials: 'same-origin', redirect: 'manual' }).catch(() => {})
    loggedIn.value = false
  }

  return { loggedIn, logout }
})
