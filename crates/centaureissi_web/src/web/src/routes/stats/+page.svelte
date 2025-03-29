<script>
    import { afterNavigate } from "$app/navigation";
    import { stats } from "$lib/api";
    import { formatBytes } from "$lib/helpers";

    let loading = false;
    let statsInfo;
    const getData = async () => {
        const statsData = await stats();
        statsInfo = statsData;
        loading =false;
    }
    afterNavigate(() => {
        loading = true;
        getData()
    })
</script>


<section class="section">
    <div class="container">
        {#if !loading}

        <div class="box">
            <h1 class="title">centaureissi stats</h1>
            {#if statsInfo}
            <p><strong>Version:</strong> {statsInfo.version}</p>
            <p><strong>Database Size:</strong> {formatBytes(statsInfo.dbSize)}</p>
            <p><strong>Mailbox Count:</strong> {statsInfo.mailboxCount} mailboxes</p>
            <p><strong>Message Count:</strong> {statsInfo.messageCount} messages</p>
            <p><strong>Blob Database Size:</strong> {formatBytes(statsInfo.blobDbSize)}</p>
            <p><strong>Blob Count:</strong> {statsInfo.blobCount} blobs</p>
            <p><strong>Search Document Count:</strong> {statsInfo.searchDocCount} documents</p>
            {:else}
            <p><strong>an issue has occured with stats gathering!</strong></p>
            {/if}
        </div>
        {:else}
        <div class="skeleton-block">

        </div>
        {/if}
    </div>
</section>