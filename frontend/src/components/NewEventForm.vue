<template>
    <div
        class="space-y-4 bg-[var(--nordic-primary-bg)] p-4 border border-[var(--nordic-border-light)] rounded-md shadow-sm w-full max-w-screen-lg">
        <h2 class="text-xl font-semibold text-[var(--nordic-text-dark)]">Create a New Event</h2>

        <!-- title -->
        <input v-model="form.title" type="text" placeholder="Event Title" required
            class="w-full px-4 py-2 border border-nordic-light rounded-md focus:outline-none focus:ring focus:ring-gray-200 text-nordic-dark" />

        <!-- description -->
        <textarea v-model="form.description" rows="4" placeholder="Event Description"
            class="w-full p-3 border border-[var(--nordic-border-light)] rounded-md 
             text-[var(--nordic-text-dark)] placeholder-[var(--nordic-text-light)] 
             bg-white
             focus:outline-none focus:ring-2 focus:ring-[var(--nordic-secondary-accent)] focus:border-[var(--nordic-secondary-accent)] resize-y"></textarea>

        <!-- date and time -->
        <input v-model="form.event_date" type="date" :min="localDate" required
            class="w-full px-4 py-2 border border-nordic-light rounded-md focus:outline-none focus:ring focus:ring-gray-200 text-nordic-dark" />
        <FlatPickr v-model="form.event_time" :config="timePickerConfig"
            class="w-full px-4 py-2 border border-nordic-light rounded-md focus:outline-none focus:ring focus:ring-gray-200 text-nordic-dark" />

        <!-- submit button -->
        <button @click="createEvent"
            :disabled="!form.title.trim() || !form.description.trim() || !form.event_date.trim() || !form.event_time.trim()"
            class="inline-flex items-center py-2 px-4 border border-[var(--nordic-border-light)] rounded-md 
             bg-[var(--nordic-primary-accent)] text-white hover:bg-[var(--nordic-secondary-accent)] 
             focus:outline-none focus:ring-2 focus:ring-[var(--nordic-secondary-accent)] transition font-medium 
             disabled:opacity-50 disabled:cursor-not-allowed">
            Create Event
        </button>

    </div>
</template>


<script setup>
import { ref, computed, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { useFormats } from '@/composables/useFormatting'
import FlatPickr from 'vue-flatpickr-component';
import 'flatpickr/dist/flatpickr.css';

const route = useRoute()
const apiUrl = import.meta.env.VITE_API_URL || '/api'
const emit = defineEmits(['event-created'])
const { localDate, localTime } = useFormats();

const timePickerConfig = reactive({
    enableTime: true,
    noCalendar: true, // Only show time picker, hide calendar
    dateFormat: "H:i", // 24-hour format (e.g., 17:30)
    time_24hr: true,
});

const form = ref({
    title: '',
    description: '',
    event_date: localDate.value,
    event_time: localTime.value,
    group_id: Number(route.params.id)
});

const event_datetime_combined = computed(() => {
    if (form.value.event_date && form.value.event_time) {
        return `${form.value.event_date}T${form.value.event_time}`;
    }
    return '';
});


const createEvent = async () => {
    try {
        const res = await fetch(`${apiUrl}/api/events/create`, {
            method: 'POST',
            body: JSON.stringify({
                group_id: form.value.group_id,
                title: form.value.title,
                description: form.value.description,
                event_datetime: event_datetime_combined.value,
            }),
            credentials: 'include',
        })

        if (res.ok) {
            const newEvent = await res.json()
            emit('event-created', newEvent)
        } else {
            alert('Failed to create event. Are you logged in?')
        }

        form.title = ''
        form.description = ''
        form.event_datetime = ''
    } catch (error) {
        console.error(error)
    }
}
</script>