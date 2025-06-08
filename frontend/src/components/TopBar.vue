<template>
    <nav class="bg-white shadow-sm w-full top-0 left-0 z-50">
        <div class="max-w-full mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex flex-wrap items-center justify-between py-4">
                <div class="flex items-center">
                    <!-- title -->
                    <router-link to="/" class="flex-shrink-0 text-2xl font-bold text-nordic-primary-accent"
                        aria-label="Home">
                        Åsocial
                    </router-link>

                    <!-- home, follows, groups, chats and events links -->
                    <div class="hidden md:flex md:ml-6 space-x-1" v-if="!isLoginPage && !isRegisterPage">
                        <router-link to="/" class="top-bar-button" :class="{ 'active': $route.path === '/' }"
                            data-title="Home" aria-label="Home">
                            <i class="fas fa-home"></i> Home
                        </router-link>
                        <router-link to="/follows" class="top-bar-button"
                            :class="{ 'active': $route.path === '/follows' }" data-title="Follows" aria-label="Follows">
                            <i class="fas fa-user-friends"></i> Follows
                        </router-link>
                        <router-link to="/groups" class="top-bar-button"
                            :class="{ 'active': $route.path === '/groups' }" data-title="Groups" aria-label="Groups">
                            <i class="fas fa-users"></i> Groups
                        </router-link>
                        <router-link to="/chats" class="top-bar-button" :class="{ 'active': $route.path === '/chats' }"
                            data-title="Chats" aria-label="Chats">
                            <i class="fas fa-comments"></i> Chats
                        </router-link>
                        <router-link to="/events" class="top-bar-button"
                            :class="{ 'active': $route.path === '/events' }" data-title="Events" aria-label="Events">
                            <i class="fas fa-calendar-alt"></i> Events
                        </router-link>
                    </div>
                </div>

                <!-- Search bar -->
                <!-- <div class="flex-1 flex justify-center px-2 lg:ml-6 lg:justify-end"
                    v-if="!isLoginPage && !isRegisterPage">
                    <div class="max-w-lg w-full lg:max-w-xs">
                        <label for="search-asocial" class="sr-only">Search</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <i class="fas fa-search text-nordic-light"></i>
                            </div>
                            <input id="search-asocial" name="search"
                                class="block w-full pl-10 pr-3 py-2 border border-nordic-light rounded-md leading-5 bg-white placeholder-nordic-light focus:outline-none focus:ring-2 focus:ring-nordic-secondary-accent focus:border-nordic-secondary-accent sm:text-sm"
                                placeholder="Search Åsocial" type="search" />
                        </div>
                    </div>
                </div> -->

                <!-- notifications and profile link -->
                <div class="hidden md:flex md:items-center md:ml-4 space-x-2" v-if="!isLoginPage && !isRegisterPage">
                    <router-link to="/notifications" class="top-bar-button relative" data-title="Notifications"
                        aria-label="Notifications">
                        <span class="sr-only">View notifications</span>
                        <i class="fas fa-bell bell"></i>
                        <span v-if="unreadCount > 0"
                            class="absolute top-0 right-0 block h-4 w-4 rounded-full ring-1 ring-white bg-red-500 text-white text-xs flex items-center justify-center">
                            {{ unreadCount }}
                        </span>
                    </router-link>

                    <router-link v-if="user" :to="`/profile/${user.id}`" class="top-bar-button"
                        :class="{ 'active': $route.path === '/profile/' + user?.id }" data-title="Your Profile"
                        aria-label="Your Profile">
                        <div class="profile-avatar">
                            <img v-if="user.avatar_url" :src="`${apiUrl}/${user.avatar_url}`" alt="User Avatar"
                                class="h-full w-full object-cover rounded-full border border-nordic-light" />
                            <i v-else class="fas fa-user"></i>
                        </div>
                        <span class="ml-2">{{ user?.first_name }}</span>
                    </router-link>

                    <button class="top-bar-button" @click="logout" data-title="Logout" aria-label="Logout">
                        <i class="fas fa-sign-out-alt"></i> Logout
                    </button>
                </div>

                <!-- hamburger menu -->
                <div class="-mr-2 flex md:hidden">
                    <button type="button"
                        class="bg-white inline-flex items-center justify-center p-2 rounded-md text-nordic-light hover:text-nordic-dark hover:bg-nordic-secondary-bg focus:outline-none focus:ring-2 focus:ring-inset focus:ring-nordic-primary-accent"
                        aria-controls="mobile-menu" :aria-expanded="isMobileMenuOpen.toString()"
                        @click="toggleMobileMenu">
                        <span class="sr-only">Open main menu</span>
                        <i :class="[isMobileMenuOpen ? 'fas fa-xmark' : 'fas fa-bars']"></i>
                    </button>
                </div>

                <!-- link to register or login -->
                <div class="hidden md:flex nav-icons" v-if="isLoginPage">
                    <router-link to="/register" class="text-nordic-light" data-title="Home"
                        aria-label="Home">Register</router-link>
                </div>
                <div class="hidden md:flex nav-icons" v-if="isRegisterPage">
                    <router-link to="/login" class="text-nordic-light" data-title="Home"
                        aria-label="Home">Login</router-link>
                </div>
            </div>
        </div>

        <!-- mobile menu -->
        <div class="md:hidden" :class="{ 'hidden': !isMobileMenuOpen }" id="mobile-menu">
            <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3" v-if="!isLoginPage && !isRegisterPage">
                <!-- home, follows, groups, chats and events links -->
                <router-link to="/" class="top-bar-button block" :class="{ 'active': $route.path === '/' }"
                    @click="toggleMobileMenu" data-title="Home" aria-label="Home">
                    <i class="fas fa-home"></i>Home
                </router-link>
                <router-link to="/follows" class="top-bar-button block"
                    :class="{ 'active': $route.path === '/follows' }" @click="toggleMobileMenu" data-title="Follows"
                    aria-label="Follows">
                    <i class="fas fa-user-friends"></i>Follows
                </router-link>
                <router-link to="/groups" class="top-bar-button block" :class="{ 'active': $route.path === '/groups' }"
                    @click="toggleMobileMenu" data-title="Groups" aria-label="Groups">
                    <i class="fas fa-users"></i>Groups
                </router-link>
                <router-link to="/chats" class="top-bar-button block" :class="{ 'active': $route.path === '/chats' }"
                    @click="toggleMobileMenu" data-title="Chats" aria-label="Chats">
                    <i class="fas fa-comments"></i>Chats
                </router-link>
                <router-link to="/events" class="top-bar-button block" :class="{ 'active': $route.path === '/events' }"
                    @click="toggleMobileMenu" data-title="Events" aria-label="Events">
                    <i class="fas fa-calendar-alt"></i>Events
                </router-link>

                <!-- notifications -->
                <router-link to="/notifications" class="top-bar-button block relative"
                    :class="{ 'active': $route.path === '/notifications' }" @click="toggleMobileMenu"
                    data-title="Notifications" aria-label="Notifications">
                    <i class="fas fa-bell bell"></i>View notifications
                    <span v-if="unreadCount > 0"
                        class="absolute top-0 left-6 block h-5 w-5 rounded-full ring-1 ring-white bg-red-500 text-white text-xs flex items-center justify-center">
                        {{ unreadCount }}
                    </span>
                </router-link>

                <!-- profile -->
                <router-link v-if="user" :to="`/profile/${user.id}`" class="top-bar-button"
                    :class="{ 'active': $route.path === '/profile/' + user?.id }" data-title="Your Profile"
                    aria-label="Your Profile">
                    <div class="profile-avatar">
                        <img v-if="user.avatar_url" :src="`${apiUrl}/${user.avatar_url}`" alt="User Avatar"
                            class="h-full w-full object-cover rounded-full border border-nordic-light" />
                        <i v-else class="fas fa-user"></i>
                    </div>
                    <span class="ml-2">{{ user?.first_name }}</span>
                </router-link>

                <!-- logout button -->
                <button
                    class="block px-3 py-2 rounded-md text-base font-medium text-nordic-light hover:text-nordic-dark hover:bg-nordic-secondary-bg w-full text-left"
                    @click="logout(); toggleMobileMenu()" data-title="Logout" aria-label="Logout">
                    <i class="fas fa-sign-out-alt"></i> Logout
                </button>
            </div>

            <!-- link to register or login -->
            <div class="nav-icons" v-if="isLoginPage">
                <router-link to="/register" class="text-nordic-light" data-title="Home"
                    aria-label="Home">Register</router-link>
            </div>
            <div class="nav-icons" v-if="isRegisterPage">
                <router-link to="/login" class="text-nordic-light" data-title="Home"
                    aria-label="Home">Login</router-link>
            </div>
        </div>
    </nav>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useUserStore } from '@/stores/user';
import { useAuth } from '@/composables/useAuth';
import { useNotificationStore } from '@/stores/notifications'; // Added

const route = useRoute();
const userStore = useUserStore();
const { user } = storeToRefs(userStore);
const { logout } = useAuth(); // Keep if used

const notificationStore = useNotificationStore(); // Added
const { unreadCount } = storeToRefs(notificationStore); // Added, 'notifications' can be used later

const isMobileMenuOpen = ref(false);
const apiUrl = import.meta.env.VITE_API_URL;
const isLoginPage = computed(() => route.path === '/login');
const isRegisterPage = computed(() => route.path === '/register');

const toggleMobileMenu = () => {
    isMobileMenuOpen.value = !isMobileMenuOpen.value;
};

onMounted(() => {
    if (userStore.isLoggedIn) { // Good practice to only fetch if logged in
        notificationStore.fetchNotifications(); // Added
    }
});

</script>

<style scoped>
/* Scoped styles specific to this component.
     Tailwind CSS classes are used directly in the template.
     The custom Nordic color classes are defined here for clarity,
     but in a real Vue project with Tailwind, you'd configure these
     colors in your tailwind.config.js file.
  */

/* Custom Nordic Colors (apply these in tailwind.config.js if using Tailwind properly) */
.text-nordic-primary-accent {
    color: #607D8B;
}

.hover\:text-nordic-secondary-accent:hover {
    color: #769FCD;
}

/* Tailwind handles hover */
.border-nordic-light {
    border-color: #CFD8DC;
}

.text-nordic-dark {
    color: #263238;
}

.text-nordic-light {
    color: #546E7A;
}

.bg-nordic-primary-accent {
    background-color: #607D8B;
}

.hover\:bg-nordic-secondary-accent:hover {
    background-color: #769FCD;
}

/* Tailwind handles hover */
.bg-nordic-secondary-bg {
    background-color: #E8EDF2;
}

.top-bar-button {
    display: flex;
    align-items: center;
    padding: 0.5rem 0.75rem;
    /* Corresponds to py-2 px-3 */
    border-radius: 0.375rem;
    /* Corresponds to rounded-md */
    transition: background-color 0.2s ease-in-out, color 0.2s ease-in-out;
    color: #546E7A;
    /* text-nordic-light */
}

.top-bar-button:hover {
    background-color: #E8EDF2;
    /* bg-nordic-secondary-bg */
    color: #263238;
    /* text-nordic-dark */
}

.top-bar-button.active {
    background-color: #E8EDF2;
    /* bg-nordic-secondary-bg */
    color: #607D8B;
    /* text-nordic-primary-accent */
    font-weight: 600;
    /* semibold */
}

.top-bar-button i {
    margin-right: 0.5rem;
    /* Corresponds to mr-2 */
}

.profile-avatar {
    width: 32px;
    /* Corresponds to w-8 */
    height: 32px;
    /* Corresponds to h-8 */
    border-radius: 9999px;
    /* Corresponds to rounded-full */
    background-color: #CFD8DC;
    /* A placeholder background */
    display: flex;
    align-items: center;
    justify-content: center;
    color: #607D8B;
    /* text-nordic-primary-accent */
}

/* Ensure Font Awesome icons are sized appropriately if needed */
.fas {
    line-height: inherit;
    /* Helps with vertical alignment */
}

.fa-user {
    margin-right: 0px !important;
}

.router-link-active.text-link,
.router-link-active.navbar-link {
    background-color: #555;
}

.bell {
        margin-right: 0px !important;
}
</style>