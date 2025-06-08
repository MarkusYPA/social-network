<template>
    <div class="chat-box">

        <div v-if="chat" ref="messagesContainer"
            class="messages-container mb-4 p-4 bg-gray-50 rounded-lg h-96 overflow-y-auto">
            <div v-if="chat.messages && chat.messages.length > 0">
                <div v-for="msg in chat.messages" :key="msg.id"
                    :class="['mb-3 p-3 rounded-lg max-w-xs', (user && msg.sender_id === user.id) ? 'ml-auto bg-nordic-primary-accent text-white' : 'bg-gray-200']">
                    <p class="text-xs font-semibold mb-1">{{ msg.sender_name }}</p>
                    <p class="break-words">{{ msg.content }}</p>
                    <span class="text-xs opacity-70 block text-right">{{ finnishTime(msg.created_at, 'medium', 'short')
                        }}</span>
                </div>
            </div>
            <div v-else class="h-full flex items-center justify-center">
                <p class="text-nordic-light italic">Start a conversation...</p>
            </div>
        </div>

        <div v-if="chat && chat.is_active" class="message-input-container">
            <form @submit.prevent="sendMessage" class="flex gap-2">
                <input id="message-input" v-model="newMessage" type="text" placeholder="Type a message..."
                    class="flex-grow p-3 bg-[var(--nordic-secondary-bg)] border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent"
                    :disabled="!isConnected" />
                <button type="submit"
                    class="px-4 py-2 bg-nordic-primary-accent text-white rounded-lg hover:bg-nordic-dark transition-colors"
                    :disabled="!newMessage.trim() || !isConnected">
                    Send
                </button>
            </form>
            <div v-if="!isConnected" class="text-red-500 text-sm mt-2">
                Not connected to chat server. Please reconnect. {{ isConnected }}
            </div>
        </div>

        <div v-else-if="chat" class="flex items-center justify-center">
            <p class="text-nordic-light italic">Follow or be followed by {{ chat.name }} to send and receive messages
            </p>
        </div>

        <div v-else class="h-64 flex items-center justify-center">
            <p class="text-nordic-light italic">Select a chat or start a new conversation</p>
        </div>
    </div>
</template>

<script setup>
// (Your script section remains the same as the last corrected version)
import { ref, watch, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useWebSocketStore } from '@/stores/websocket'
import { storeToRefs } from 'pinia';
import { useUserStore } from '@/stores/user';
import { useFormats } from '@/composables/useFormatting'

const props = defineProps({
    chat: {
        type: Object,
        default: null
    },
    groupString: {      // put "group" in front of message type when in group chat
        type: String,
        default: ''
    }
})

const route = useRoute()
const newMessage = ref('')
const websocketStore = useWebSocketStore()
const { isConnected } = storeToRefs(websocketStore)
const userStore = useUserStore()
const { user } = storeToRefs(userStore)
const messagesContainer = ref(null);
const { finnishTime } = useFormats();

// Note: Ensure this watch block is aligned with the latest ChatBox.vue changes for deduplication.
// It should primarily handle echoes for the senderconst { finnishTime } = useFormats();

watch(() => websocketStore.message, (newMsg) => {
    if (newMsg && newMsg.type === `${props.groupString}chat_message` && props.chat) {

        // If this is an echo of the current user's message, perform deduplication
        if (Number(newMsg.from) === Number(user.value.id) && Number(newMsg.receiver_id) === Number(props.chat.user_id)) {
            const msgIdSent = newMsg.id; // Assuming server echoes back the client-generated ID
            const messageAlreadyAdded = props.chat.messages?.some(msg => msg.id === msgIdSent);

            if (messageAlreadyAdded) {
                // console.log("Duplicate self-sent message from echo detected in ChatBox, skipping.");
                return;
            }
            // If for some reason, the sent message wasn't added locally,
            // or if the server provides a canonical ID you want to use for the existing message,
            // you might have logic here to find and update/replace that local message.
            // For now, if messageAlreadyAdded is true, we simply return.
        }

        // private chat messages handled in ChatsView.vue, group chat messages in ChatBox.vue

        if (!props.chat.messages) {
            props.chat.messages = [];
        }

        // Check if message belongs to this groupchat
        if (newMsg.type === 'groupchat_message' && props.groupString === 'group' && newMsg.receiver_id === route.params.id) {
            props.chat.messages.push({
                id: `${newMsg.from}_${Date.now()}`, // Unique ID for the message
                content: newMsg.content,
                sender_name: newMsg.from_name, // Use from_name from websocket message
                sender_id: Number(newMsg.from),
                created_at: new Date(newMsg.timestamp || Date.now()) // Use timestamp from WS or current
            });
        }
    }
})

function sendMessage() {
    if (!newMessage.value.trim() || !props.chat || !isConnected.value) return;

    const msgId = `${user.value.id}_${props.chat.user_id}_${Date.now()}`;

    const message = {
        id: msgId,
        created_at: new Date(),
        content: newMessage.value,
        sender_name: user.value.first_name,
        sender_id: user.value.id
    };

    if (props.chat.messages) {
        props.chat.messages.push(message);
    } else {
        props.chat.messages = [message];
    }

    const wsMessage = {
        type: `${props.groupString}chat_message`,
        from: user.value.id,
        from_name: user.value.first_name,
        receiver_id: props.chat.user_id || '0', // Use the user ID from the chat object, means group id when group chat
        content: newMessage.value
    }

    websocketStore.send(wsMessage);

    newMessage.value = '';
}

watch(
    () => props.chat?.messages,
    () => {
        nextTick(() => {
            if (messagesContainer.value) {
                messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
            }
        })
    },
    { deep: true }
)

onMounted(() => {
    nextTick(() => {
        if (messagesContainer.value) {
            messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
        }
    })
})
</script>