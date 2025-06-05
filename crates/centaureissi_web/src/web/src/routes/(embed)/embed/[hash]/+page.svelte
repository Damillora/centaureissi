<script lang="ts">
    const { data } = $props();
    import { afterNavigate } from "$app/navigation";
    import { getMessageByType, getMessageDetail } from "$lib/api";
    import { secondsInDay } from "date-fns/constants";
    import { onMount } from "svelte";
    import { format, formatDistanceToNow } from "date-fns";

    let content = $state("");

    const loadMessage = async (message_type: string) => {
        const response = await getMessageByType(data.hash, message_type);

        content = response.content;
    };

    onMount(async () => {
        loadMessage("html");
    });
</script>

{@html content}
