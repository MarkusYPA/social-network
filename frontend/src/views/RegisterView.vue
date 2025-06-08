<template>
    <div class="register">
        <TopBar />

        <div class="px-4 mt-20 flex justify-center">
            <div class="w-full max-w-md bg-white p-8 rounded-lg shadow-md border border-nordic-light">
                <h2 class="text-2xl font-bold text-center text-nordic-light mb-6">Register</h2>

                <form @submit.prevent="register" class="flex flex-col gap-4">
                    <input v-model="form.email" type="email" placeholder="Email" required autocomplete="email"
                        class="p-3 border border-nordic-light rounded-md focus:outline-none  focus:ring-2 focus:ring-nordic-primary-accent" />

                    <input v-model="form.password" type="password" placeholder="Password" required
                        autocomplete="new-password"
                        class="p-3 border border-nordic-light rounded-md focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent" />

                    <input v-model="form.firstName" type="text" placeholder="First Name" required
                        class="p-3 border border-nordic-light rounded-md focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent" />

                    <input v-model="form.lastName" type="text" placeholder="Last Name" required
                        class="p-3 border border-nordic-light rounded-md focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent" />

                    <div class="flex flex-col">
                        <label for="dob" class="ml-3 mb-1 text-nordic-light">Date of Birth</label>
                        <input v-model="form.dob" type="date"  :max="localDate" required
                            class="p-3 border border-nordic-light rounded-md focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent" />
                    </div>

                    <input v-model="form.nickname" type="text" placeholder="Nickname (Optional)"
                        class="p-3 border border-nordic-light rounded-md focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent" />

                    <textarea v-model="form.about" placeholder="About Me (Optional)" rows="3"
                        class="w-full px-4 py-2 border border-nordic-light rounded-md focus:outline-none focus:ring focus:ring-gray-200 resize-none text-nordic-dark"></textarea>

                    <div class="space-y-2">
                        <label class="block text-sm font-medium text-nordic-light">Avatar/Image (Optional):</label>
                        <input type="file" @change="handleFileUpload" accept="image/*"
                            class="block w-full text-sm text-nordic-light file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-gray-100 file:text-nordic-light hover:file:bg-gray-200" />

                        <div v-if="form.avatarUrl" class="mt-2">
                            <p class="text-sm text-nordic-light">Current Avatar:</p>
                            <img :src="`${apiUrl}/${form.avatarUrl}`" alt="Avatar"
                                class="h-24 w-24 object-cover rounded-md border border-nordic-light" />
                            <button type="button" @click="deleteAvatar"
                                class="mt-2 text-sm text-red-600 hover:underline text-nordic-primary-accent hover:text-nordic-secondary-accent">Delete
                                Avatar</button>
                        </div>
                        <!-- <div v-else>
                    <span class="block text-sm font-medium text-nordic-light">No current avatar</span>
                </div> -->
                    </div>

                    <button type="submit"
                        class="bg-nordic-primary-accent text-white font-medium py-2 rounded-md hover:bg-nordic-secondary-accent transition">Register</button>

                    <p v-if="error" class="text-red-600 text-center">{{ error }}</p>
                </form>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import TopBar from '@/components/TopBar.vue';
import { useImageProcessor } from '@/composables/useImageProcessor';
import { useFormats } from '@/composables/useFormatting'

const form = ref({
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    dob: '',
    nickname: '',
    about: '',
    avatar: null,
});

const error = ref(null);
const router = useRouter();
const apiUrl = import.meta.env.VITE_API_URL
const { localDate } = useFormats();

const handleFileUpload = (event) => {
    form.value.avatar = event.target.files[0];
};

const register = async () => {
    const { processAvatarImage } = useImageProcessor();

    try {
        const payload = {
            email: form.value.email,
            password: form.value.password,
            firstName: form.value.firstName,
            lastName: form.value.lastName,
            dob: form.value.dob,
            nickname: form.value.nickname || null,
            about: form.value.about || null
        };

        const formData = new FormData();
        for (const [key, value] of Object.entries(payload)) {
            formData.append(key, value);
        }
        if (form.value.avatar) {
            const processedAvatar = await processAvatarImage(form.value.avatar);
            formData.append('avatar', processedAvatar);
        }

        const response = await fetch(`${apiUrl}/api/register`, {
            method: 'POST',
            body: formData,
            credentials: 'include'
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Registration failed');
        }

        // Handle success (e.g., redirect to login)
        router.push('/login');
    } catch (err) {
        error.value = err.message || 'An error occurred';
    }
};
</script>
