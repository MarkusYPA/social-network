
<template>
    <div>
        <TopBar />

        <div class="px-4 mt-20 flex justify-center">
            <div class="w-full max-w-md bg-white p-8 rounded-lg shadow-md border border-nordic-light">
                <h2 class="text-2xl font-bold text-center text-nordic-light mb-6">Login</h2>

                <form @submit.prevent="login" class="flex flex-col gap-4">
                    <input v-model="email" type="email" placeholder="Email" required autocomplete="email"
                        class="p-3 border border-nordic-light rounded-md focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent" />
                    <input v-model="password" type="password" placeholder="Password" required
                        autocomplete="current-password"
                        class="p-3 border border-nordic-light rounded-md focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent" />
                    <button type="submit"
                        class="bg-nordic-primary-accent text-white font-medium py-2 rounded-md hover:bg-nordic-secondary-accent transition">Login</button>
                    <p v-if="err">{{ err }}</p>
                </form>

                <!-- Add this divider and GitHub login button -->
                <div class="flex items-center my-6">
                    <div class="flex-grow border-t border-nordic-light"></div>
                    <span class="mx-4 text-nordic-light">OR</span>
                    <div class="flex-grow border-t border-nordic-light"></div>
                </div>
                
                <button @click="loginWithGitHub"
                    class="w-full bg-gray-800 text-white font-medium py-2 rounded-md hover:bg-gray-700 transition flex items-center justify-center gap-2">
                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                        <path fill-rule="evenodd" d="M10 0C4.477 0 0 4.484 0 10.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0110 4.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.203 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.942.359.31.678.921.678 1.856 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0020 10.017C20 4.484 15.522 0 10 0z" clip-rule="evenodd" />
                    </svg>
                    Continue with GitHub
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import TopBar from '@/components/TopBar.vue'
import { useErrorStore } from '@/stores/error'
import { useWebSocketStore } from '@/stores/websocket'
import { storeToRefs } from 'pinia';

const email = ref('')
const password = ref('')
const err = ref('')
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const wsUrl = import.meta.env.VITE_WS_URL
const wsStore = useWebSocketStore()
const { user } = storeToRefs(userStore)


async function login() {

    try {
        const res = await fetch(`${apiUrl}/api/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ email: email.value, password: password.value }),
        })
        if (res.ok) {
            const data = await res.json()

            //console.log("data to store at login:", data.user)
            userStore.setUser(data.user)

            /*             wsStore.connect(`${wsUrl}/ws`)   // main.js should take care of connecting
                        watch(() => wsStore.message, (newMessage) => {
                            if (newMessage?.type === 'welcome') {
                                console.log('WebSocket authenticated successfully')
                            }
                        }) */


            /*             type WSMessage struct {
                Type    string `json:"type"`
                From    string `json:"from"`
                To      string `json:"receiver_id,omitempty"`
                Content string `json:"content,omitempty"`
            } */


            // Testing connection
            const nuMsg = {
                type: "ping",
                from: user.value.id,
                from_name: user.value.first_name + ' ' + user.value.last_name          
            }
            wsStore.send(nuMsg)

            // Navigate to what the user wanted or home,
            let redirectTo = route.query.redirect || '/'
            router.push(redirectTo)
        } else {
            const msg = await res.text()
            err.value = msg || 'Login failed'
        }

    } catch (err) {
        errorStore.setError('Unexpected Error', 'Something went wrong while logging in.')
        router.push('/error')
        return
    }

}

async function loginWithGitHub() {
    window.location.href = `${apiUrl}/api/auth/github`;
}
</script>
