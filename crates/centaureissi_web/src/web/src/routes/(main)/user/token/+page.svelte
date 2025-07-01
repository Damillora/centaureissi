<script>
    import {
        getUserProfile,
        getToken,
        updateToken,
        updateUserProfile,
    } from "$lib/api";
    import UserActionsPanel from "$lib/components/panels/UserActionsPanel.svelte";
    import { onMount } from "svelte";

    let loading = $state(false);
    let auth_token = $state("");
    let tokencopy = $state();

    const submitForm = async () => {
        loading = true;
        const token = await getToken();
        auth_token = token;
        loading = false;
    };
    const copyToken = () => {
        tokencopy.select();
        tokencopy.setSelectionRange(0, 99999);
        document.execCommand("copy");
    };
</script>

<section class="section">
    <div class="container">
        <div class="columns">
            <div class="column is-one-third">
                <UserActionsPanel />
            </div>
            <div class="column is-two-thirds">
                <div class="box">
                    <h1 class="title">Token</h1>
                    <form onsubmit={submitForm}>
                        <div class="field">
                            <p>
                                Click to generate a 24-hour token for importer.
                            </p>
                        </div>
                        <div class="field">
                            <input
                                class="input"
                                type="text"
                                readonly
                                placeholder="Token"
                                value={auth_token}
                                class:is-skeleton={loading}
                                onclick={copyToken}
                                bind:this={tokencopy}
                            />
                        </div>
                        <div class="field">
                            <button
                                class="button is-primary is-fullwidth is-outlined"
                                type="submit">Generate Token</button
                            >
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</section>
