<template>
    <div class="group-view">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <div v-if="group">
                    <!-- name and description -->
                    <h2 class="text-2xl font-bold mb-4">{{ group.title }}</h2>
                    <p>{{ group.description }}</p>

                    <!-- members only information -->
                    <div v-if="membershipStatus === 'accepted' || membershipStatus === 'admin'">
                        <br>
                        <!-- chat button -->
                        <button @click="toggleChat()"
                            class="mb-4 px-4 py-2 bg-nordic-primary-accent hover:bg-nordic-secondary-accent text-white rounded transition">
                            {{ chatOpen ? 'Close Chat' : 'Open Chat' }}
                        </button>
                        <!-- events -->
                        <EventList :events="events" small class="my-4" />
                        <!-- members -->
                        <MembersList :members="members" small class="my-4" />
                    </div>

                    <GroupReqNoticesForAdmin v-if="membershipStatus === 'admin'"
                        @update-members='getMembers(route.params.id)' />
                </div>
            </template>

            <template #main>

                <!-- button to join / leave / delete -->
                <button @click="prepareGroupAction" class="mb-4 px-4 py-2 rounded transition" :class="groupButtonClass">
                    {{
                        membershipStatus === '' ? 'Request to Join' :
                            membershipStatus === 'pending' ? 'Cancel Request' :
                                membershipStatus === 'accepted' ? 'Leave Group' :
                                    membershipStatus === 'declined' ? 'Request to Join' :
                                        membershipStatus === 'admin' ? 'Delete Group' :
                                            ''
                    }}
                </button>

                <!-- title image -->
                <div class="relative">
                    <img :src="`${apiUrl}/data/default/groupdefault01.jpg`" alt="Page Image"
                        class="w-full max-w-screen-lg mb-4 h-40 object-cover rounded" />
                    <div class="absolute inset-0 flex items-center justify-center">
                        <h1 v-if="group && group.title" class="text-white text-5xl font-extrabold text-center">
                            {{ group.title }}
                        </h1>
                    </div>
                </div>

                <!-- members only content -->
                <div v-if="membershipStatus === 'accepted' || membershipStatus === 'admin'">

                    <div v-if="!chatOpen">
                        <!-- buttons for new post, invite user and new event -->
                        <span class="flex flex-wrap gap-4">
                            <button @click="showPostForm = !showPostForm; showInviteForm = false; showEventForm = false"
                                class="mb-4 px-4 py-2 bg-nordic-primary-accent hover:bg-nordic-secondary-accent text-white rounded transition">
                                {{ showPostForm ? 'Cancel Post' : 'New Post' }}
                            </button>
                            <button
                                @click="showInviteForm = !showInviteForm; showPostForm = false; showEventForm = false"
                                class="mb-4 px-4 py-2 bg-nordic-primary-accent hover:bg-nordic-secondary-accent text-white rounded transition">
                                {{ showInviteForm ? 'Close Invitation Form' : 'Invite Users' }}
                            </button>
                            <button
                                @click="showEventForm = !showEventForm; showPostForm = false; showInviteForm = false"
                                class="mb-4 px-4 py-2 bg-nordic-primary-accent hover:bg-nordic-secondary-accent text-white rounded transition">
                                {{ showEventForm ? 'Close Event Form' : 'New Event' }}
                            </button>
                        </span>

                        <!-- forms to create new post, invite user or new event -->
                        <NewGroupPostForm v-if="showPostForm" :group_id="Number(route.params.id)"
                            @post-submitted="handlePostSubmitted" class="mb-8" />
                        <InviteUsers v-if="showInviteForm" :members="members" class="mb-8" />
                        <NewEventForm v-if="showEventForm" @event-created="handleEventCreated" class="mb-8" />

                        <!-- group posts -->
                        <PostsList :posts="posts" />
                    </div>

                    <!-- chat -->
                    <div v-else>
                        <ChatBox :chat="groupChat" group-string="group" />
                        <!-- literal string passed to prop without : -->
                    </div>

                </div>
            </template>
        </TwoColumnLayout>

        <ConfirmDialog :visible="showLeaveConfirmation" title="Leave Group"
            message="Are you sure you want to leave this group? This action cannot be undone."
            @confirm="handleGroupAction" @cancel="showLeaveConfirmation = false" />

        <ConfirmDialog :visible="showDeleteConfirmation" title="Delete Group"
            message="Are you sure you want to delete this group? This action cannot be undone."
            @confirm="handleGroupAction" @cancel="showDeleteConfirmation = false" />
    </div>
</template>

<script setup>
import { onMounted, ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useErrorStore } from '@/stores/error'
import TopBar from '@/components/TopBar.vue'
import PostsList from '@/components/PostsList.vue'
import EventList from '@/components/EventList.vue'
import MembersList from '@/components/MembersList.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import NewGroupPostForm from '@/components/NewGroupPostForm.vue'
import GroupReqNoticesForAdmin from '@/components/GroupReqNoticesForAdmin.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import InviteUsers from '@/components/InviteUsers.vue'
import ChatBox from '@/components/ChatBox.vue'
import NewEventForm from '@/components/NewEventForm.vue'

const route = useRoute()
const router = useRouter()
const errorStore = useErrorStore()
const apiUrl = import.meta.env.VITE_API_URL
const group = ref(null)
const posts = ref([])
const members = ref([])
const events = ref([])
const showPostForm = ref(false)
const showInviteForm = ref(false)
const showEventForm = ref(false)
const showLeaveConfirmation = ref(false)
const showDeleteConfirmation = ref(false)
const membershipStatus = ref('')

const chatOpen = ref(false)
const groupChat = ref(null)

const handlePostSubmitted = (newPost) => {
    posts.value.unshift(newPost)
}

const handleEventCreated = (newEvent) => {
    //console.log("new event:", newEvent)
    events.value.push(newEvent)
    showEventForm.value = false
}

const groupButtonClass = computed(() => {
    if (membershipStatus.value === '' || membershipStatus.value === 'declined') {
        return 'bg-nordic-primary-accent hover:bg-nordic-secondary-accent text-white';
    }
    if (membershipStatus.value === 'pending') {
        return 'bg-nordic-secondary-bg hover:bg-nordic-secondary-accent text-nordic-light';
    }
    if (membershipStatus.value === 'accepted') {
        //return 'bg-nordic-text-light hover:bg-nordic-primary-accent text-black';  // doesn't work for some reason
        return 'bg-nordic-primary-accent hover:bg-nordic-secondary-accent text-white';
    }
    if (membershipStatus.value === 'admin') {
        return 'bg-nordic-primary-accent hover:bg-nordic-secondary-accent text-white';
    }
    return '';
});

function toggleChat() {
    chatOpen.value = !chatOpen.value
    getChat(route.params.id)
    //console.log("Chat open:", chatOpen.value)
}

async function getChat(groupId) {        // Fetch and filter posts
    try {
        const postsRes = await fetch(`${apiUrl}/api/group/chat/messages/${groupId}`, {     //
            credentials: 'include'
        })

        if (!postsRes.ok) {
            throw new Error(`Failed to fetch posts: ${postsRes.status}`)
        }

        groupChat.value = await postsRes.json()
        //console.log("chat call succeeded")
    } catch (error) {
        console.log("error fetching group chat:", error)
        errorStore.setError('Error', 'Something went wrong while loading group chat data.')
        router.push('/error')
        return
    }
}

function prepareGroupAction() {
    if (!showLeaveConfirmation.value && membershipStatus.value === 'accepted') {
        showLeaveConfirmation.value = true
        return
    }
    if (!showDeleteConfirmation.value && membershipStatus.value === 'admin') {
        showDeleteConfirmation.value = true
        return
    }
    handleGroupAction()
}

async function handleGroupAction() {
    let action = ''
    if (membershipStatus.value === '' || membershipStatus.value === 'declined') action = 'request'
    else if (membershipStatus.value === 'accepted') action = 'leave'
    else if (membershipStatus.value === 'pending') action = 'cancel'
    else if (membershipStatus.value === 'admin') action = 'delete'

    try {
        const res = await fetch(`${apiUrl}/api/group/join`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({
                target_id: group.value.id,
                action: action
            })
        })

        if (!res.ok) throw new Error('Failed to update group status')

        // Get membership status again
        if (action != 'delete') {
            const followRes = await fetch(`${apiUrl}/api/group/${group.value.id}`, {     //
                credentials: 'include'
            })
            if (!followRes.ok) {
                throw new Error(`Failed to fetch follow info in handleGroupAction: ${followRes.status}`)
            }
            const resp = await followRes.json()
            membershipStatus.value = resp.membership
        }
    } catch (err) {
        console.error(err)
    }

    showLeaveConfirmation.value = false
    showDeleteConfirmation.value = false

    if (action == 'delete') {
        router.push('/groups')
    }
}

async function getGroup(groupId) {
    try {
        // Fetch user info
        const groupRes = await fetch(`${apiUrl}/api/group/${groupId}`, {
            credentials: 'include'
        })

        if (groupRes.status === 401) { // Session is invalid — logout and redirect
            logout()
            router.push('/login')
            return
        }

        if (groupRes.status === 404) {
            errorStore.setError('Group Not Found', `Group with ID ${userId} does not exist.`)
            router.push('/error')
            return
        }

        if (groupRes.status === 400) {
            errorStore.setError('Bad request', `Failed to get group with ID ${userId}.`)
            router.push('/error')
            return
        }

        if (!groupRes.ok) {
            // Generic error
            throw new Error(`Failed to fetch group: ${groupRes.status}`)
        }

        const resp = await groupRes.json()
        group.value = resp.group
        membershipStatus.value = resp.membership

    } catch (err) {
        console.log("error fetching group:", err)
        errorStore.setError('Error', 'Failed to get group data.')
        router.push('/error')
        return
    }
}

async function getPosts(groupId) {        // Fetch and filter posts
    try {
        const postsRes = await fetch(`${apiUrl}/api/group/posts/${groupId}`, {     //
            credentials: 'include'
        })

        if (!postsRes.ok) {
            throw new Error(`Failed to fetch posts: ${postsRes.status}`)
        }
        const groupPosts = await postsRes.json()
        if (groupPosts) posts.value = groupPosts
    } catch (error) {
        console.log("error fetching group posts:", error)
        errorStore.setError('Error', 'Something went wrong while loading group posts data.')
        router.push('/error')
        return
    }
}

async function getMembers(groupId) {        // Fetch and filter posts
    try {
        const membsRes = await fetch(`${apiUrl}/api/group/members/${groupId}`, {     //
            credentials: 'include'
        })

        if (!membsRes.ok) {
            throw new Error(`Failed to fetch posts: ${membsRes.status}`)
        }
        const groupMembs = await membsRes.json()
        if (groupMembs) members.value = groupMembs
    } catch (error) {
        console.log("error fetching group members:", error)
        errorStore.setError('Error', 'Something went wrong while loading group members data.')
        router.push('/error')
        return
    }
}

async function getEvents(groupId) {        // Fetch and filter posts
    try {
        const eventsRes = await fetch(`${apiUrl}/api/group/events/${groupId}`, {     //
            credentials: 'include'
        })

        if (!eventsRes.ok) {
            throw new Error(`Failed to fetch posts: ${eventsRes.status}`)
        }
        const groupEvs = await eventsRes.json()
        if (groupEvs) events.value = groupEvs
    } catch (error) {
        console.log("error fetching group events:", error)
        errorStore.setError('Error', 'Something went wrong while loading group events data.')
        router.push('/error')
        return
    }
}

onMounted(() => {
    getGroup(route.params.id)
    getPosts(route.params.id)
    getMembers(route.params.id)
    getEvents(route.params.id)
    getChat(route.params.id)
})

watch(() => membershipStatus, () => {
    getGroup(route.params.id)
    getPosts(route.params.id)
    getMembers(route.params.id)
    getEvents(route.params.id)
})
</script>
