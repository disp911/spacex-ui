import { ref } from 'vue'

// One breakpoint: below it the panel switches to its phone layout
// (drawer navigation, full-screen sheets, cards instead of tables).
const query = window.matchMedia('(max-width: 760px)')
export const isMobile = ref(query.matches)
query.addEventListener('change', (e) => (isMobile.value = e.matches))
