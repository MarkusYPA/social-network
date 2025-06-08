import { computed } from 'vue';

export const useFormats = () => {
    const pad = n => n.toString().padStart(2, '0');
    const now = new Date();

    const currentTime = computed(() => {
        const now = new Date();
        now.setSeconds(0, 0); // Remove seconds and milliseconds
        return now.toISOString().slice(0, 16);
    })

    const finnishTime = (isoString, dStyle, mStyle) => {    // "medium", "short"
        const date = new Date(isoString);
        return date.toLocaleString("fi-FI", {
            dateStyle: dStyle,
            timeStyle: mStyle,
        }).replace("klo ", "");
    }


    const localDate = computed(() => {
        return `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())}`;
    });

    const localTime = computed(() => {
        return `${pad(now.getHours())}:${pad(now.getMinutes())}`;
    });

    const datetimeLocal = computed(() => {
        return `${localDate.value}T${localTime.value}`;
    });

    return {
        currentTime,
        finnishTime,
        localDate,
        localTime,
        datetimeLocal
    }
}