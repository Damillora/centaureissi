<script lang="ts">
    import { afterNavigate, goto } from "$app/navigation";
    import { search } from "$lib/api";
    import AuthRequired from "$lib/components/checks/AuthRequired.svelte";
    import { paginate } from "$lib/simple-pagination";
    import { page as currentPage } from "$app/state";
    import MailItem from "$lib/components/ui/MailItem.svelte";

    let url = $derived(currentPage.url);

    let query = $state("");
    let page = $state(1);
    let totalPages = $state(1);
    let pagination: string[] = $state([]);
    let messages: any[] = $state([]);
    let messageCount = $state(0);

    let loading = $state(false);

    const getData = async () => {
        if (query == "") {
            loading = false;
            return;
        }
        const data = await search({ q: query, page: page });
        if (data.items) {
            messages = data.items;
            totalPages = data.totalPages;
            messageCount = data.count;
            pagination = paginate(page, totalPages);
        }

        loading = false;
    };

    afterNavigate(() => {
        loading = true;
        query = url.searchParams.get("q");
        messages = [];
        page = 1;
        getData();
    });

    const onSearch = (e) => {
        e.preventDefault();
        goto(`/search?q=${query}`);
    };

    const changePage = (i) => {
        if (i >= 1 && i <= totalPages) {
            page = i;
            getData();
        }
    };
</script>

<AuthRequired />

<section class="section">
    <div class="container">
        <div class="block">
            <div class="column is-full">
                <form onsubmit={onSearch}>
                    <div class="field has-addons">
                        <div class="control is-expanded">
                            <div class="control" id="search">
                                <input
                                    class="input"
                                    type="text"
                                    bind:value={query}
                                />
                            </div>
                        </div>
                        <div class="control">
                            <button type="submit" class="button is-primary">
                                Search
                            </button>
                        </div>
                    </div>
                </form>
            </div>
        </div>
        {#if !loading}
            <div class="block">
                <div class="column is-full">
                    <div class="block">
                        {#each messages as message, i (message.id)}
                            <MailItem {message} />
                        {/each}
                    </div>
                    <nav class="pagination is-centered" aria-label="pagination">
                        <a
                            href={null}
                            onclick={() => changePage(page - 1)}
                            class="pagination-previous"
                            class:is-disabled={page == 1}>Previous</a
                        >
                        <a
                            href={null}
                            onclick={() => changePage(page + 1)}
                            class="pagination-next"
                            class:is-disabled={page == totalPages}>Next</a
                        >
                        <ul class="pagination-list">
                            {#each pagination as pageEntry}
                                {#if pageEntry == "..."}
                                    <li>
                                        <span class="pagination-ellipsis"
                                            >&hellip;</span
                                        >
                                    </li>
                                {:else}
                                    <li>
                                        <a
                                            href={null}
                                            onclick={() =>
                                                changePage(pageEntry)}
                                            class="pagination-link"
                                            class:is-current={page == pageEntry}
                                            aria-label="Goto page {pageEntry}"
                                            >{pageEntry}</a
                                        >
                                    </li>
                                {/if}
                            {/each}
                        </ul>
                    </nav>
                </div>
            </div>
        {:else}
            <div class="skeleton-block"></div>
        {/if}
    </div>
</section>
