<template>
    <div v-if="results && results.length > 0" class="mb-8">
        <h3 class="text-xl font-semibold text-[var(--nordic-text-dark)] mb-3">Search Results</h3>
        <ul class="space-y-2">
            <li v-for="user in results" :key="user.user.id"
                class="text-[var(--nordic-text-light)] hover:text-[var(--nordic-primary-accent)] transition-colors duration-150  break-all">

                <!-- invite button -->
                <button @click="inviteUser(user.user.id)"
                    :disabled="user.membership === 'invited' || user.membership === 'accepted'"
                    class="w-40 bg-[var(--nordic-primary-accent)] text-white px-2 py-1 mr-4 rounded hover:bg-[var(--nordic-secondary-accent)] disabled:opacity-50 disabled:cursor-not-allowed">
                    <span v-if="user.membership === ''">Invite</span>
                    <span v-else-if="user.membership === 'invited'">Invited</span>
                    <span v-else>Unavailable</span>
                </button>

                <!-- link to profile -->
                <RouterLink :to="`/profile/${user.user.id}`">
                    {{ user.user.first_name }} {{ user.user.last_name }}
                    <span v-if="user.user.username"> - {{ user.user.username }}</span>
                </RouterLink>
            </li>
        </ul>
    </div>
    <p v-else-if="searchInitiated" class="text-[var(--nordic-text-light)] italic mb-6">No users found.</p>
</template>

<script setup>
import { useRoute } from 'vue-router'

const props = defineProps({
    results: Array,
    searchInitiated: Boolean
})

const route = useRoute()
const apiUrl = import.meta.env.VITE_API_URL

const inviteUser = async (userId) => {
    const groupId = route.params.id
    try {
        const response = await fetch(`${apiUrl}/api/group/invite`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({
                user_id: userId,
                group_id: groupId
            })
        })

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`)
        }

        // Update local copy user membership status
        const userIndex = props.results.findIndex(user => user.user.id === userId)
        if (userIndex !== -1) {
            props.results[userIndex].membership = 'invited'
        }

    } catch (error) {
        console.error('Failed to invite user:', error)
    }
}
</script>
