<template>
    <div
        class="post-card flex flex-col items-start gap-2 mb-4 p-4 bg-[var(--nordic-primary-bg)] border border-[var(--nordic-border-light)] rounded-md shadow-sm">

        <img v-if="post.image_path" :src="`${apiUrl}/${post.image_path}`" alt=""
            class="w-full rounded-md border border-[var(--nordic-border-light)]" />

        <p class="post-content text-nordic-dark text-base break-all">
            {{ post.content }}
        </p>

        <small class="post-date flex items-center text-sm">
            <RouterLink :to="`/profile/${post.user_id}`"
                class="post-user flex items-center mr-1 hover:underline text-[var(--nordic-primary-accent)] break-all">
                <div v-if="post.avatar_url"
                    class="post-user-avatar w-6 h-6 rounded-full overflow-hidden mr-1 border border-[var(--nordic-border-light)]">
                    <img :src="`${apiUrl}/${post.avatar_url}`" alt="User Avatar" class="w-full h-full object-cover" />
                </div>
                {{ post.username }}
            </RouterLink>
            <span v-if="post.group_id">
                in
                <RouterLink :to="`/group/${post.group_id}`"
                    class="text-[var(--nordic-primary-accent)] hover:underline mr-1 break-all">
                    {{ post.group_name }}
                </RouterLink>
            </span>
            on {{ formattedPostDate }}
        </small>

        <div class="flex flex-row w-full justify-between items-center">
            <div class="flex items-center">
                <button @click="toggleComments" class="group relative flex items-center mr-4 p-1 rounded-md transition-colors duration-200"
                    :disabled="commentCount === 0"
                    :class="{ 'text-[var(--nordic-border-light)] cursor-default': commentCount === 0, 'text-[var(--nordic-secondary-accent)] hover:bg-blue-100': commentCount > 0 }">
                    <font-awesome-icon :icon="showComments ? 'eye-slash' : 'comments'" class="text-xl" />
                    <span v-if="commentCount > 0" class="ml-1 text-sm font-semibold">{{ commentCount }}</span>
                    <span
                        class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-2 py-1 bg-gray-700 text-white text-xs rounded-md opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap pointer-events-none">
                        {{ commentLabel }}
                    </span>
                </button>

                <button @click="showNewCommentForm" class="group relative flex items-center p-1 rounded-md text-[var(--nordic-secondary-accent)] hover:bg-blue-100 transition-colors duration-200">
                    <font-awesome-icon :icon="newComment ? 'comment-slash' : 'comment-medical'" class="text-xl" />
                    <span
                        class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-2 py-1 bg-gray-700 text-white text-xs rounded-md opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap pointer-events-none">
                        {{ newComment ? 'Hide form' : 'Add comment' }}
                    </span>
                </button>
            </div>

            <span class="text-xs break-normal text-gray-400">{{ displayPrivacy() }}</span>
        </div>

        <div v-if="newComment" class="space-y-4 bg-white p-4 border rounded-md shadow-sm w-full max-w-screen-lg">
            <h2 class="text-xl font-semibold text-gray-800">Create a new comment</h2>
            <textarea v-model="content" rows="4" placeholder="Enter your comment" class="w-full p-3 border border-gray-300 rounded-md text-gray-700 placeholder-gray-400
                focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-y"></textarea>
            <div class="space-y-2">
                <label class="block text-sm font-medium text-[var(--nordic-text-light)]">Image (Optional):</label>
                <input type="file" @change="handleFileUpload" accept=".jpg, .jpeg, .gif, .png, .webp" ref="imageInput" class="block w-full text-sm text-[var(--nordic-text-light)]
                    file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0
                    file:text-sm file:font-semibold
                    file:bg-[var(--nordic-secondary-bg)] file:text-[var(--nordic-text-dark)]
                    hover:file:bg-[var(--nordic-border-light)]" />
            </div>
            <button @click="submitComment" :disabled="!content.trim()" class="inline-flex items-center py-2 px-4 border border-nordic-light rounded-md
                bg-nordic-primary-accent text-white hover:bg-nordic-secondary-accent focus:outline-none
                focus:ring-2 focus:ring-nordic-secondary-accent transition font-medium
                disabled:opacity-50 disabled:cursor-not-allowed">
                submit comment
            </button>
        </div>

        <div v-if="showComments" class="mt-2">
            <div v-if="loadingComments">Loading comments...</div>
            <div v-else-if="commentCount === 0">No comments yet.</div>
            <ul v-else>
                <li v-for="comment in comments" :key="comment.id" class="mt-1 border-t pt-1">
                    <img v-if="comment.image_path" :src="`${apiUrl}/${comment.image_path}`" alt=""
                        class="w-full rounded-md border border-[var(--nordic-border-light)]" />
                    <p class="post-content text-[var(--nordic-text-dark)] text-base break-all">
                        {{ comment.content }}
                    </p>
                    <small class="post-date flex items-center text-sm text-[var(--nordic-text-light)]">
                        <RouterLink :to="`/profile/${post.user_id}`"
                            class="post-user flex items-center mr-1 hover:underline text-[var(--nordic-text-dark)] break-all">
                            <div v-if="comment.user.avatar_url"
                                class="post-user-avatar w-6 h-6 rounded-full overflow-hidden mr-1 border border-[var(--nordic-border-light)]">
                                <img :src="`${apiUrl}/${comment.user.avatar_url}`" alt="User Avatar"
                                    class="w-full h-full object-cover" />
                            </div>
                            {{ comment.user.first_name }} {{ comment.user.last_name }}
                        </RouterLink>
                        on {{ finnishTime(comment.created_at, 'medium', 'medium') }}
                    </small>
                </li>
            </ul>
        </div>

    </div>
</template>


<script setup>
import { computed, ref, watch, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useImageProcessor } from '@/composables/useImageProcessor';
import { useFormats } from '@/composables/useFormatting'

const { post } = defineProps({
    post: {
        type: Object,
        required: true
    }
})

const comments = ref([]);
const content = ref('');
const showComments = ref(false);
const newComment = ref(false);
const loadingComments = ref(false);
const apiUrl = import.meta.env.VITE_API_URL || '/api';
const image = ref(null);
const imageInput = ref(null); // Add ref for file input
const commentCount = ref(post.numberOfComments || 0);  // Track comment count separately
const { finnishTime } = useFormats();

const handleFileUpload = (event) => {
    image.value = event.target.files[0];
};

const displayPrivacy = () => {
    if (post.privacy === 'public') return 'public post'
    if (post.privacy === 'almost_private') return 'almost private post'
    if (post.privacy === 'private') return 'private post'
    return 'group post'
}


// Update comment count when comments are loaded
watch(comments, (newComments) => {
    commentCount.value = newComments.length;
});

const commentLabel = computed(() => {
    if (commentCount.value === 0) {
        return 'No comments yet';
    }
    return showComments.value
        ? `Hide Comments`
        : `Show Comments`;
});

const formattedPostDate = computed(() => finnishTime(post.created_at, 'medium', 'short'))

const showNewCommentForm = async () => {
    newComment.value = !newComment.value;
};

const toggleComments = async () => {
    if (commentCount.value === 0) return;

    showComments.value = !showComments.value;
    if (showComments.value) {
        await loadComments();
    }
};

const loadComments = async () => {
    loadingComments.value = true;
    try {
        const res = await fetch(`${apiUrl}/api/comments/show?post_id=${post.id}&type=${post.postType}`, {
            credentials: 'include',
        });
        comments.value = await res.json();
        // Update the count after loading
        commentCount.value = comments.value?.length || 0;
        //console.log("Image is: ", comments.)
    } catch (error) {
        comments.value = [];
        console.error("Failed to load comments", error);
    } finally {
        loadingComments.value = false;
    }
};

const submitComment = async () => {
    const { processPostImage } = useImageProcessor();

    try {

        const payload = {
            content: content.value,
            type: post.postType,
        };

        const formData = new FormData();
        for (const [key, value] of Object.entries(payload)) {
            formData.append(key, value);
        }

        if (image.value) {
            const processedImg = await processPostImage(image.value);
            formData.append('image', processedImg);
        }

        const res = await fetch(`${apiUrl}/api/comments/create?post_id=${post.id}`, {
            method: 'POST',
            body: formData,
            credentials: 'include',
        })
        const result = await res.json();
        if (result.success) {
            content.value = '';
            image.value = null;
            if (imageInput.value) {
                imageInput.value.value = ''; // Clear file input
            }
            showNewCommentForm()
            // Increment comment count
            commentCount.value++;
            // Reload comments if they're shown
            if (showComments.value) {
                await loadComments();
            }
        } else {
            alert('Failed to comment. Are you logged in?')
        }
    } catch (error) {
        console.error(error)
    }
};

/* onMounted(() => {
    if (post.postType !== 'group' && post.postType === 'regular') {
        console.log("This regular post:", post)
    }
}) */

</script>

<style scoped>

.post-content {
    font-size: 1.1rem;
}

.post-date {
    color: #666;
    font-size: 0.85rem;
}

.post-user {
    text-decoration: none;
    font-weight: bold;
}

.post-user:hover {
    text-decoration: underline;
}
</style>