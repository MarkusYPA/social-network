<template>
    <li class="notification-item">
        <small class="post-date flex items-center">
            {{ formatDate(notification.created_at) }}
        </small>
        <div class="notification-content">
            <span class="notification-text">
                <template v-if="notification.type === 'follow_request'">
                    <router-link :to="'/profile/' + notification.sender_id" class="font-bold">{{
                        notification.sender_name }}</router-link> wants to follow you
                </template>
                <template v-else-if="notification.type === 'group_invitation'">
                    <router-link :to="'/profile/' + notification.sender_id" class="font-bold">{{
                        notification.sender_name }}</router-link> invited you to join
                    <router-link :to="'/group/' + notification.group_id" class="font-bold">{{
                        notification.group_title }}</router-link>
                </template>
                <template v-else-if="notification.type === 'group_join_request'">
                    <router-link :to="'/profile/' + notification.sender_id" class="font-bold">{{
                        notification.sender_name }}</router-link> wants to join
                    <router-link :to="'/group/' + notification.group_id" class="font-bold">{{
                        notification.group_title }}</router-link>
                </template>
                <template v-else-if="notification.type === 'event_creation'">
                    New event: <router-link :to="'/events/' + notification.event_id" class="font-bold">{{
                        notification.event_title }}</router-link>
                </template>
                <template v-else>
                    {{ notification.content }}
                </template>
            </span>

            <!-- Close Button -->
            <button class="close-button" @click="closeNotification(notification.id)" v-if="!notification.is_read">
                &times;
            </button>
        </div>

        <!-- Accept/Decline Buttons -->
        <div class="action-buttons" v-if="notification.pending">
            <button class="accept-button" @click="acceptAction(notification.id)">
                {{ notification.type === 'event_creation' ? 'Going' : 'Accept' }}
            </button>
            <button class="decline-button" @click="declineAction(notification.id)">
                {{ notification.type === 'event_creation' ? 'Not going' : 'Decline' }}
            </button>
        </div>
    </li>
</template>

<script setup>

function formatDate(dateString) {
    const date = new Date(dateString)
    return date.toLocaleString("ru-RU")
}

// Accept props
defineProps({
    notification: {
        type: Object,
        required: true
    }
})

// Emit events to the parent component
const emit = defineEmits(['close', 'accept', 'decline'])

function closeNotification(id) {
    emit('close', id)
}

function acceptAction(id) {
    emit('accept', id)
}

function declineAction(id) {
    emit('decline', id)
}
</script>

<style scoped>
.notification-item {
    position: relative;
    padding: 1rem;
    margin-bottom: 0.5rem;
    background-color: #f9f9f9;
    border-left: 5px solid #888;
    border-radius: 4px;
}

.notification-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.notification-text {
    flex: 1;
    color: var(--nordic-text-dark);
}

/* Close Button (X) */
.close-button {
    position: absolute;
    top: 0.5rem;
    right: 0.75rem;
    background: none;
    border: none;
    font-size: 1.5rem;
    color: var(--nordic-text-light);
    cursor: pointer;
    transition: color 0.2s ease-in-out;
}

.close-button:hover {
    color: var(--nordic-primary-accent);
}

/* Invitation Action Buttons */
.action-buttons {
    margin-top: 0.75rem;
    display: flex;
    gap: 0.5rem;
}

.accept-button,
.decline-button {
    padding: 0.4rem 1rem;
    font-size: 0.9rem;
    border-radius: 20px;
    border: none;
    cursor: pointer;
    transition: background-color 0.2s ease-in-out;
}

.accept-button {
    background-color: var(--nordic-primary-accent);
    color: white;
}

.accept-button:hover {
    background-color: var(--nordic-secondary-accent);
}

.decline-button {
    background-color: var(--nordic-secondary-bg);
    color: var(--nordic-text-dark);
}

.decline-button:hover {
    background-color: #dbe4ec;
}
</style>