<script>
    import { formatDistanceToNow } from "date-fns";

    let { message } = $props();

    let viewMore = $state(false);

    const toggleViewMore = () => {
        viewMore = !viewMore;
    };
</script>

<article class="media">
    <!-- <figure class="media-left">
    </figure> -->
    <div class="media-content">
        <div class="content">
            <p>
                <small>From: {message.from}</small>
                <br />
                <a href={`/message/${message.hash}`}
                    ><strong>{message.subject}</strong></a
                >
                <small
                    >{formatDistanceToNow(message.date, {
                        addSuffix: true,
                    })}</small
                >
                <br />
                <a href={"#"} onclick={toggleViewMore}><small>Details</small></a
                >
                {#if viewMore}
                    <br />
                    <small>To: {message.to}</small>
                    <br />
                    <small>Cc: {message.cc}</small>
                    <br />
                    <small>Bcc: {message.bcc}</small>
                {/if}
            </p>
        </div>
    </div>
    <!-- <div class="media-right">
    </div> -->
</article>
