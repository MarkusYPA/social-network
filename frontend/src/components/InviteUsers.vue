<template>
    <div class="space-y-4 bg-[var(--nordic-primary-bg)] p-4 border rounded-md shadow-sm w-full max-w-screen-lg">
        <h2 class="text-xl font-semibold text-[var(--nordic-text-dark)]">Invite Users</h2>

        <SearchBoxInvite @results="handleResults" />
        <SearchResultsInvite :results="searchResults" :searchInitiated="searchInitiated" />
    </div>
</template>

<script setup>
import { ref } from 'vue'
import SearchBoxInvite from './SearchBoxInvite.vue'
import SearchResultsInvite from './SearchResultsInvite.vue'

//const emit = defineEmits(['post-submitted'])
const searchResults = ref([])
const searchInitiated = ref(false)

const { members } = defineProps({
    members: {
        type: Array,
        required: true
    }
});

function handleResults(results) {
    searchInitiated.value = true
    searchResults.value = results.filter(invitableUser => {
        const isMember = members.some(member => member.id === invitableUser.user.id)
        return !isMember
    });
}

</script>