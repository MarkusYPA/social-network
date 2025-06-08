<template>
    <div class="notifications-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h3 class="text-lg font-semibold">Notifications</h3>
                <button @click="showCurrent = true" :class="{ active: showCurrent }">
                    Current ({{ unreadCount }})
                </button>
                <button @click="showCurrent = false" :class="{ active: !showCurrent }">
                    <div v-if="notifications && notifications.length > 0">
                        Old ({{ notifications.length - unreadCount }})
                    </div>
                    <div v-else>
                        Old (0)
                    </div>
                </button>
            </template>

            <template #main>
                <h2 class="text-2xl font-bold mb-4">{{ showCurrent ? 'Current' : 'Old' }} Notifications</h2>
                <NotificationsList :notifications="filteredNotifications" @close="handleClose" @accept="handleAccept"
                    @decline="handleDecline" />
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import NotificationsList from '@/components/NotificationsList.vue'
import { useErrorStore } from '@/stores/error'
import { useAuth } from '@/composables/useAuth'
import { useNotificationStore } from '@/stores/notifications';
import { storeToRefs } from 'pinia';
import { useUserStore } from '@/stores/user';   // necessary?

const showCurrent = ref(true)
const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const router = useRouter()
const { logout } = useAuth()
const notificationStore = useNotificationStore()
const { notifications, unreadCount } = storeToRefs(notificationStore)
const userStore = useUserStore()
const { user } = storeToRefs(userStore)

const filteredNotifications = computed(() =>
    (notifications.value || []).filter(n => n.is_read !== showCurrent.value)
);

// Old fetchNotifications and readNotification functions are removed.

async function approveFollowRequest(id, action) {
    try {
        const res = await fetch(`${apiUrl}/api/follow/requests/${id}/${action}`, {
            credentials: 'include'
        })

        if (res.status === 401) {
            logout();
            router.push('/login');
            return;
        }

        if (!res.ok) throw new Error(`Failed to accept/decline follow request: ${res.status}`)

    } catch (err) {
        errorStore.setError('Error', `Error while accepting/declining follow request`)
        router.push('/error')
    }
}

async function approveGroupRequest(groupID, senderID, action) {
    //console.log(action, "to group request for group", groupID, "from user", senderID)

    try {
        const res = await fetch(`${apiUrl}/api/group/requests/${action}`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                group_id: groupID,
                requester_id: senderID
            })
        })

        if (res.status === 401) {
            logout();
            router.push('/login');
            return;
        }

        if (!res.ok) throw new Error(`Failed to accept/decline group request: ${res.status}`)

    } catch (err) {
        errorStore.setError('Error', `Error while accepting/declining group request`)
        router.push('/error')
    }
}


async function approveGroupInvite(groupInviteID, action) {
    try {
        const res = await fetch(`${apiUrl}/api/group/invite/${groupInviteID}/${action}`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
        })

        if (res.status === 401) {
            logout();
            router.push('/login');
            return;
        }

        if (!res.ok) throw new Error(`Failed to accept/decline group invitation: ${res.status}`)

    } catch (err) {
        errorStore.setError('Error', `Error while accepting/declining group invitation`)
        router.push('/error')
    }
}

async function answerToEvent(eventID, action) {
    if (action === 'accepted') action = 'going'
    if (action === 'declined') action = 'not_going'

    try {
        const res = await fetch(`${apiUrl}/api/events/respond`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                event_id: eventID,
                user_id: user.value.id,
                response: action
            })
        })

        if (res.status === 401) {
            logout();
            router.push('/login');
            return;
        }

        if (!res.ok) throw new Error(`Failed to accept/decline group invitation: ${res.status}`)

    } catch (err) {
        errorStore.setError('Error', `Error while accepting/declining group invitation`)
        router.push('/error')
    }
}

onMounted(() => {
    // Fetch notifications only if the store is empty, assuming TopBar might have loaded them
    if (!notifications.value || notifications.value?.length === 0) {
        notificationStore.fetchNotifications();
    }
});

async function handleClose(id) { // Made async to align with potential async operations in store
    const n = notifications.value.find(n => n.id === id);
    if (!n) return;
    await notificationStore.markAsRead(id); // Use store action
}

async function handleAccept(id) {
    const n = notifications.value.find(n => n.id === id); // Use store's notifications
    if (!n) return;

    await notificationStore.markAsRead(id); // Mark as read via store

    try {
        switch (n.type) {
            case 'follow_request':
                await approveFollowRequest(n.follow_req_id, 'accept'); // Existing API call
                break;
            case 'group_join_request':
                await approveGroupRequest(n.group_id, n.sender_id, 'accepted'); // Existing API call
                break;
            case 'group_invitation':
                await approveGroupInvite(n.group_invite_id, 'accepted'); // Existing API call
                break;
            case 'event_creation':
                await answerToEvent(n.event_id, 'accepted'); // Existing API call
                break;
        }
        // After successful action, refresh the list from the store to get latest state
        await notificationStore.fetchNotifications();
    } catch (err) {
        // errorStore.setError is already called by approveFollowRequest etc.
        // or add specific error handling here if needed
        console.error(`Error handling accept for notification ${id}:`, err);
    }
}

async function handleDecline(id) {
    const n = notifications.value.find(n => n.id === id); // Use store's notifications
    if (!n) return;

    await notificationStore.markAsRead(id); // Mark as read via store

    try {
        switch (n.type) {
            case 'follow_request':
                await approveFollowRequest(n.follow_req_id, 'decline'); // Existing API call
                break;
            case 'group_join_request':
                await approveGroupRequest(n.group_id, n.sender_id, 'declined'); // Existing API call
                break;
            case 'group_invitation':
                await approveGroupInvite(n.group_invite_id, 'declined'); // Existing API call
                break;
            case 'event_creation':
                await answerToEvent(n.event_id, 'declined'); // Existing API call
                break;
        }
        // After successful action, refresh the list from the store to get latest state
        await notificationStore.fetchNotifications();
    } catch (err) {
        // errorStore.setError is already called by approveFollowRequest etc.
        // or add specific error handling here if needed
        console.error(`Error handling decline for notification ${id}:`, err);
    }
}

</script>

<style scoped>
.notifications-page {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}

button {
    display: block;
    margin: 0.5rem 0;
    padding: 0.5rem 1rem;
    background-color: #eee;
    border: none;
    cursor: pointer;
    border-radius: 5px;
}

button.active {
    background-color: #ccc;
    font-weight: bold;
}
</style>
