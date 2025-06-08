import './assets/main.css'

import { createApp, watchEffect } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import { useUserStore } from './stores/user'
import { useWebSocketStore } from './stores/websocket'
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faComments, faCommentMedical, faEyeSlash, faCommentSlash } from '@fortawesome/free-solid-svg-icons' // Import necessary icons

library.add(faComments, faCommentMedical, faEyeSlash, faCommentSlash) // Add them to the library
const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)
app.component('font-awesome-icon', FontAwesomeIcon)
const userStore = useUserStore(pinia)
const websocketStore = useWebSocketStore(pinia)
const apiUrl = import.meta.env.VITE_API_URL || '/api'

watchEffect(() => {
    if (userStore.isLoggedIn) {
        websocketStore.initWebSocket()
    } else {
        websocketStore.disconnect()
    }
})

try {
    const res = await fetch(`${apiUrl}/api/me`, {
        credentials: 'include'
    });

    if (res.status !== 200) {
        throw new Error('Not logged in');
    }

    const user = await res.json();
    userStore.setUser(user);
    app.mount('#app'); // mount app after successful credentials check
} catch (error) {
    // If not logged in and not already on login page, redirect to login
    if (router.currentRoute.value.path !== '/login' && router.currentRoute.value.path !== '/register') {
        if (router.currentRoute.value.path === '/error') {  // don't redirect to error page
            router.replace('/login')
        } else {
            router.replace({ path: '/login', query: { redirect: router.currentRoute.value.fullPath } });
        }
    }
    app.mount('#app');
}
