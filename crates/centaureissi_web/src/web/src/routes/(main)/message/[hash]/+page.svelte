<script lang="ts">
    const { data } = $props();
    import { afterNavigate } from "$app/navigation";
    import { getMessageByType, getMessageDetail } from "$lib/api";
    import { secondsInDay } from "date-fns/constants";
    import { onMount } from "svelte";
    import { format, formatDistanceToNow } from "date-fns";

    let loading = $state(true);
    let loading_content = $state(true);

    let message_data: any = $state({
        subject: "",
    });
    let message_types: string[] = $state([]);
    let selected_type = $state("");
    let content = $state("");

    const getDetail = async () => {
        const response = await getMessageDetail(data.hash);

        message_data = { ...response };

        let newtypes: string[] = [];
        if (message_data.is_html_mail) {
            newtypes.push("html");
        }
        if (message_data.is_text_mail) {
            newtypes.push("text");
        }
        newtypes.push("raw");

        message_types = newtypes;
        selected_type = newtypes[0];
        await loadMessage(selected_type);

        loading = false;
    };
    const onChangeType = async (message_type: string) => {
        selected_type = message_type;
        await loadMessage(selected_type);
    };
    const loadMessage = async (message_type: string) => {
        loading_content = true;
        const response = await getMessageByType(data.hash, message_type);

        content = response.content;
        loading_content = false;
    };

    onMount(async () => {
        loading = true;
        message_data = {};
        await getDetail();
    });
</script>

<section class="section">
    <div class="container">
        <div class="columns">
            <div class="column is-one-third">
                {#if !loading}
                    <div class="panel is-primary">
                        <p class="panel-heading">Message</p>
                        <div class="panel-block column">
                            <div class="row">
                                <strong>Date:</strong>
                            </div>
                            <div class="row">
                                <time
                                    title={format(
                                        message_data.date,
                                        "dd MMMM yyyy HH:mm:ss",
                                    )}
                                    datetime={message_data.date}
                                    >{formatDistanceToNow(message_data.date, {
                                        addSuffix: true,
                                    })}</time
                                >
                            </div>
                        </div>
                        <div class="panel-block column">
                            <div class="row">
                                <strong>From:</strong>
                            </div>
                            <div class="row">
                                <p>{message_data.from}</p>
                            </div>
                        </div>
                        <div class="panel-block column">
                            <div class="row">
                                <strong>To:</strong>
                            </div>
                            <div class="row">
                                <p>{message_data.to}</p>
                            </div>
                        </div>
                        {#if message_data.cc}
                            <div class="panel-block column">
                                <div class="row">
                                    <strong>CC:</strong>
                                </div>
                                <div class="row">
                                    <p>{message_data.cc}</p>
                                </div>
                            </div>
                        {/if}
                        {#if message_data.bcc}
                            <div class="panel-block column">
                                <div class="row">
                                    <strong>BCC:</strong>
                                </div>
                                <div class="row">
                                    <p>{message_data.bcc}</p>
                                </div>
                            </div>
                        {/if}
                    </div>
                {:else}
                    <div class="skeleton-block"></div>
                {/if}
            </div>
            <div class="column is-two-thirds">
                {#if !loading}
                    <div class="block">
                        {#if !loading_content}
                            <div class="panel is-primary">
                                <div class="panel-heading">
                                    <h1 class="title">
                                        {message_data.subject}
                                    </h1>
                                </div>
                                <p class="panel-tabs">
                                    {#each message_types as message_type}
                                        <a
                                            href={"#"}
                                            class:is-active={selected_type ===
                                                message_type}
                                            onclick={() =>
                                                onChangeType(message_type)}
                                        >
                                            {message_type}
                                        </a>
                                    {/each}
                                </p>
                                <div class="panel-row block">
                                    {#if selected_type === "html"}
                                        <figure class="image is-16by9">
                                            <iframe
                                                class="has-ratio"
                                                title={message_data.subject}
                                                src={`/embed/${data.hash}`}
                                            >
                                            </iframe>
                                        </figure>
                                    {:else if selected_type === "text"}
                                        <div class="content">
                                            <pre>{content}</pre>
                                        </div>
                                    {:else}
                                        <div class="content">
                                            <pre>{content}</pre>
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        {:else}
                            <div class="skeleton-block"></div>
                        {/if}
                    </div>
                {:else}
                    <div class="skeleton-block"></div>
                {/if}
            </div>
        </div>
    </div>
</section>
