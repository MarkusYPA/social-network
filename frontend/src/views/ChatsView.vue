<template>
    <div class="chats-page-wrapper">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>

                <div class="mb-8">

                    <div class="flex justify-between items-center mb-3">
                        <h3 class="text-xl font-semibold text-nordic-dark">Chats</h3>
                        <span v-if="isConnected"
                            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                            Online
                        </span>
                        <span v-else
                            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
                            Offline
                        </span>
                    </div>

                    <ul v-if="chats.length > 0" class="space-y-2">
                        <li v-for="chat in chats" :key="chat.user_id" @click="select(chat)" :class="[ // Use chat.user_id for key
                            'cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 p-2 rounded-lg',
                            selectedChat && selectedChat.user_id === chat.user_id ? 'font-semibold text-nordic-primary-accent bg-gray-100' : ''
                        ]">
                            <span>{{ chat.name }}</span>
                            <div class="text-sm text-gray-500 truncate" v-if="getLastMessage(chat)">
                                {{ getLastMessage(chat).sender_name }}: {{ getLastMessage(chat).content }}
                            </div>
                            <span class="text-xs text-gray-500" v-if="getLastMessage(chat)">
                                {{ finnishTime(getLastMessage(chat).created_at, 'medium', 'short') }}
                            </span>

                        </li>
                    </ul>
                    <p v-else class="italic text-nordic-light">No chats yet.</p>
                </div>


                <h4 class="text-xl font-semibold text-nordic-dark mb-3">Connections</h4>

                <div class="mb-4">
                    <h4 class="text-lg font-medium text-nordic-dark">Following</h4>
                    <ul v-if="followedUsers.length > 0" class="space-y-2">
                        <li v-for="user in followedUsers" :key="user.id">
                            <button @click="startChat(user)"
                                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                                {{ user.first_name }} {{ user.last_name }}
                            </button>
                        </li>
                    </ul>
                    <p v-else class="italic text-nordic-light">You follow no one yet.</p>
                </div>

                <div>
                    <h4 class="text-lg font-medium text-nordic-dark">Followers</h4>
                    <ul v-if="followers.length > 0" class="space-y-2">
                        <li v-for="user in followers" :key="user.id">
                            <button @click="startChat(user)"
                                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                                {{ user.first_name }} {{ user.last_name }}
                            </button>
                        </li>
                    </ul>
                    <p v-else class="italic text-nordic-light">No followers yet.</p>
                </div>
            </template>

            <template #main>
                <div v-if="selectedChat" class="mb-6">
                    <div class="flex justify-between items-center">
                        <h2 class="text-3xl font-bold text-nordic-dark">{{ selectedChat.name }}</h2>
                        <button v-if="!isConnected" @click="connectWebSocket"
                            class="px-4 py-2 bg-nordic-primary-accent text-white rounded-lg hover:bg-nordic-dark transition-colors">
                            Reconnect
                        </button>
                    </div>
                </div>
                <ChatBox v-if="selectedChat" :chat="selectedChat" />
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useWebSocketStore } from '@/stores/websocket'
import TopBar from '@/components/TopBar.vue'
import ChatBox from '@/components/ChatBox.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia';
import { useAuth } from '@/composables/useAuth'
import { useErrorStore } from '@/stores/error'
import { useRouter } from 'vue-router'
import { useFormats } from '@/composables/useFormatting'

const websocketStore = useWebSocketStore()
const isConnected = computed(() => websocketStore.isConnected)
const wsUrl = import.meta.env.VITE_WS_URL
const userStore = useUserStore()
const { user } = storeToRefs(userStore)
const followedUsers = ref([])
const followers = ref([])
const { logout } = useAuth()
const router = useRouter()
const errorStore = useErrorStore()
const apiUrl = import.meta.env.VITE_API_URL
const { finnishTime } = useFormats();

const selectedChat = ref(null)
const chats = ref([])


onMounted(async () => {
    fetchFollowData()
    testConnection()
    await getAllUserChats()

    // Set default selected_chat
    if (chats.value.length > 0) {
        selectedChat.value = chats.value[0]
    }
})

function connectWebSocket() {
    websocketStore.connect(`${wsUrl}/ws`) // Corrected connect call
}

function testConnection() {
    const nuMsg = {
        type: "ping",
        from: user.value.id,
    }
    websocketStore.send(nuMsg)
}

// listen for new webSocket chat_messages:
watch(() => websocketStore.message, async (message) => { // Made async to await fetch
    if (!message || message.type !== 'chat_message') return;

    // Check if the message is for the current user
    // IMPORTANT: Ensure message.receiver_id and user.value.id are consistently strings or numbers.
    // If from backend, they might be numbers. user.value.id from Pinia might be number.
    // Using loose equality (==) helps, but explicit conversion is safer.
    if (Number(message.receiver_id) !== Number(user.value.id)) {
        console.log("Message not for current user:", message.receiver_id, user.value.id);
        return;
    }

    // Check if current user's chats already include one with the same sender
    const existingChat = chats.value.find(c => Number(c.user_id) === Number(message.from)); // Use message.from as sender_id

    // If not existing, create a new one
    if (!existingChat) {
        try {
            // Fetch sender info based on sender_id
            const response = await fetch(`${apiUrl}/api/users/${message.from}`, { credentials: 'include' });
            if (!response.ok) {
                throw new Error('Failed to fetch user info');
            }
            const senderData = await response.json();
        
            // Create a new chat object with correct structure
            const newChat = {
                is_active: true, // Assuming new chats are active by default
                name: `${senderData.first_name} ${senderData.last_name}`,
                user_id: Number(message.from), // Use user_id for consistency
                messages: [] // Start with an empty messages array
            };

            // Add the new chat to the chats list
            chats.value.push(newChat);

            // Add the *received message* to the newly created chat's messages
            newChat.messages.push({
                id: `${message.from}_${Date.now()}`, // Unique ID for the message
                content: message.content,
                sender_name: message.from_name, // Use from_name from websocket message
                sender_id: Number(message.from), // Use from as sender_id
                created_at: new Date(message.timestamp || Date.now()) // Use timestamp from WS or current
            });

            // If this new chat is the currently selected one (e.g., if it's the only chat
            // or if the user was waiting for this specific message), select it.
            // This is a decision point: Do you want to auto-select new chats?
            // For now, let's auto-select if no chat is selected.
            if (!selectedChat.value) {
                selectedChat.value = newChat;
            }

        } catch (error) {
            console.error('Error fetching user info or creating new chat:', error);
            // Optionally, push an error message to the user
            errorStore.setError('Chat Error', 'Could not start new chat with user.');
        }
    } else {
        // If chat exists, add the message to its messages array
        console.log("Chat exists, adding message to it.");
        // Ensure messages array exists
        if (!existingChat.messages) {
            existingChat.messages = [];
        }

        // Check for duplicates to prevent adding the same message multiple times
        const messageExists = existingChat.messages.some(msg =>
            Number(msg.sender_id) === Number(message.from) &&
            msg.content === message.content &&
            Math.abs(new Date(msg.created_at).getTime() - new Date(message.timestamp || Date.now()).getTime()) < 5000 // Deduplicate within 5 seconds
        );

        if (!messageExists) {
            existingChat.messages.push({
                id: `${message.from}_${Date.now()}`, // Unique ID for the message
                content: message.content,
                sender_name: message.from_name, // Use from_name from websocket message
                sender_id: Number(message.from), // Use from as sender_id
                created_at: new Date(message.timestamp || Date.now()) // Use timestamp from WS or current
            });
        }
    }
});


function select(chat) {
    selectedChat.value = chat
}

function startChat(userToChatWith) { // Renamed 'user' param to avoid conflict with Pinia 'user'
    const existingChat = chats.value.find(chat => Number(chat.user_id) === Number(userToChatWith.id))

    if (existingChat) {
        select(existingChat)
    } else {
        const newChat = {
            is_active: true,
            name: `${userToChatWith.first_name} ${userToChatWith.last_name}`,
            user_id: Number(userToChatWith.id), // Ensure consistent type
            messages: [] // New chat starts with empty messages
        }
        chats.value.push(newChat)
        selectedChat.value = newChat
    }
}

function getLastMessage(chat) {
    if (!chat.messages || chat.messages.length === 0) return null
    return chat.messages[chat.messages.length - 1]
}

async function fetchFollowData() {
    try {

        const [res1, res2] = await Promise.all([
            fetch(`${apiUrl}/api/followed/0`, { credentials: 'include' }),
            fetch(`${apiUrl}/api/followers/0`, { credentials: 'include' })
        ])

        if (res1.status === 401 || res2.status === 401) {
            logout()
            router.push('/login')
            return
        }

        if (!res1.ok || !res2.ok) throw new Error('Failed to fetch follow data')

        const [followedJson, followersJson] = await Promise.all([
            await res1.json(),
            await res2.json()
        ])

        if (followedJson) followedUsers.value = followedJson
        if (followersJson) followers.value = followersJson

    } catch (err) {
        errorStore.setError('Error', 'Failed to load follow data.')
        router.push('/error')
        console.log('error fetching in chatsVue')
    }
}

async function getAllUserChats() {
    try {
        const chatResp = await fetch(`${apiUrl}/api/chat/messages`, { credentials: 'include' })

        if (chatResp.status === 401) {
            logout()
            router.push('/login')
            return
        }

        if (!chatResp.ok) {
            const errorData = await chatResp.json().catch(() => ({ message: 'Failed to fetch posts and parse error.' })); // Corrected 'res' to 'chatResp'
            // isLoading.value = false // Assuming no isLoading in this component
            throw new Error(errorData.message || `HTTP error ${chatResp.status}`)
        }

        const fetchedChats = await chatResp.json();
        // Ensure fetched chats have a messages array (even if empty) and proper types
        chats.value = fetchedChats.map(chat => ({
            ...chat,
            user_id: Number(chat.user_id), // Ensure user_id is a number
            messages: chat.messages ? chat.messages.map(msg => ({
                id: msg.id || `${msg.sender_id}_${msg.created_at || Date.now()}`, // Ensure ID for existing messages
                content: msg.content,
                sender_name: msg.sender_name,
                sender_id: Number(msg.sender_id), // Ensure sender_id is a number
                created_at: msg.created_at // Assuming created_at is valid date string/object
            })) : [] // Ensure messages is always an array
        }));

    } catch (error) {
        errorStore.setError('Error Loading Chats', error.message || 'An unexpected error occurred while trying to load chats. Please try again later.');
        router.push('/error')
        return
    }
}


</script>

<style scoped>
.chats-page-wrapper {
    min-height: 100vh;
}
</style>