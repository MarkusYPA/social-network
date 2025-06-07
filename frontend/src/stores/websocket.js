import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useUserStore } from '../stores/user'
import { useNotificationStore } from './notifications' // Added import

export const useWebSocketStore = defineStore('websocket', () => {
    const socket = ref(null)
    const isConnected = ref(false)
    const message = ref(null)
    const reconnectAttempts = ref(0)
    const maxReconnectAttempts = 5
    const reconnectTimeout = ref(null)
    const userStore = useUserStore()

    function connect(url) {
        if (socket.value) return // already connected or connecting

        clearTimeout(reconnectTimeout.value)
        socket.value = new WebSocket(url)

        socket.value.onopen = () => {
            isConnected.value = true
            reconnectAttempts.value = 0
            console.log('WebSocket connected')

            // Send a ping to keep the connection alive
            send({
                type: 'ping',
                content: 'ping',
                sender_id: userStore.userId || '0',
                receiver_id: '0'
            })
        }

        socket.value.onmessage = (event) => {
            try {
                // Parse the incoming message
                message.value = JSON.parse(event.data)
                //console.log('WebSocket message received:', message.value); // Generic log

                // Handle 'new_notification'
                if (message.value.type === 'new_notification') {
                    //console.log('Received new_notification message content:', message.value.content);
                    const notificationStore = useNotificationStore(); // Get store instance
                    let newNotificationData;
                    try {
                        if (typeof message.value.content === 'string') {
                            newNotificationData = JSON.parse(message.value.content);
                        } else {
                            // If content is already an object
                            newNotificationData = message.value.content; 
                        }
                    } catch (parseError) {
                        console.error('Error parsing new_notification content:', parseError, "Raw content:", message.value.content);
                        return; // Skip if content is not valid JSON
                    }

                    if (newNotificationData) {
                        notificationStore.addNotification(newNotificationData);
                    }
                } 
                // Handle ping response (existing logic)
                else if (message.value.type === 'pong') {
                    console.log('Received pong from server:', message.value.content);
                }
                // Handle 'notification_deleted'
                else if (message.value.type === 'notification_deleted') {
                    console.log('Received notification_deleted message content:', message.value.content);
                    const notificationStore = useNotificationStore(); // Get store instance

                    let deletedNotificationData;
                    try {
                        // Assuming message.value.content is a JSON string like '{"id":"123"}'
                        if (typeof message.value.content === 'string') {
                            deletedNotificationData = JSON.parse(message.value.content);
                        } else {
                            // If content is already an object (less likely based on backend plan but good to handle)
                            deletedNotificationData = message.value.content;
                        }
                    } catch (parseError) {
                        console.error('Error parsing notification_deleted content:', parseError, "Raw content:", message.value.content);
                        return; // Skip if content is not valid JSON
                    }

                    if (deletedNotificationData && deletedNotificationData.id) {
                        // The removeNotification function in notifications.js expects a number or a string that can be coerced to a number.
                        // Ensure the ID is passed correctly.
                        const notificationIdToRemove = Number(deletedNotificationData.id);
                        if (!isNaN(notificationIdToRemove)) {
                            notificationStore.removeNotification(notificationIdToRemove);
                        } else {
                            console.error('Invalid notification ID received for deletion:', deletedNotificationData.id);
                        }
                    } else {
                        console.warn('Notification ID missing in notification_deleted message content:', deletedNotificationData);
                    }
                }

            } catch (error) {
                console.error('Error parsing message envelope:', error, "Raw data:", event.data);
                // message.value = event.data; // Avoid setting message.value to raw on envelope parse error
            }
        }

        socket.value.onclose = () => {
            isConnected.value = false
            socket.value = null
            console.log('WebSocket closed')

            // Try to reconnect if not maxed out
            if (reconnectAttempts.value < maxReconnectAttempts) {
                reconnectAttempts.value++
                const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.value), 30000) // Exponential backoff with 30s max
                console.log(`Attempting to reconnect in ${delay}ms (attempt ${reconnectAttempts.value})`)

                reconnectTimeout.value = setTimeout(() => {
                    connect(url)
                }, delay)
            }
        }

        socket.value.onerror = (error) => {
            console.error('WebSocket error:', error)
        }
    }

    function send(data) {
        if (!isConnected.value || !socket.value) {
            console.warn('Cannot send message: WebSocket not connected')
            return false
        }

        try {
            // Ensure all numeric values are sent as strings to avoid parsing issues in Go
            const processedData = Object.entries(data).reduce((acc, [key, value]) => {
                // Convert any numbers to strings
                if (typeof value === 'number') {
                    acc[key] = value.toString()
                } else {
                    acc[key] = value
                }
                return acc
            }, {})

            socket.value.send(JSON.stringify(processedData))
            //console.log("sending processed data:", JSON.stringify(processedData))
            return true
        } catch (error) {
            console.error('Error sending message:', error)
            return false
        }
    }

    function disconnect() {
        if (socket.value && isConnected.value) {
            socket.value.close()
        }

        // Clear any pending reconnect attempts
        clearTimeout(reconnectTimeout.value)
        reconnectAttempts.value = maxReconnectAttempts // Prevent further reconnects
    }

    function initWebSocket() {
        const userStore = useUserStore()
        if (userStore.isLoggedIn) {
            const websocketUrl = import.meta.env.VITE_WEBSOCKET_URL || 'ws://localhost:8080/ws'
            connect(websocketUrl)
        }
    }

    return { connect, send, disconnect, isConnected, message, initWebSocket }
})