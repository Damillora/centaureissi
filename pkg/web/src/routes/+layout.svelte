<script lang="ts">
    import "../app.scss";
    import Navbar from "$lib/components/ui/Navbar.svelte";
    import Footer from "$lib/components/ui/Footer.svelte";
    interface Props {
        children?: import('svelte').Snippet;
    }

    let { children }: Props = $props();
</script>

<svelte:head>
    <title>centaureissi</title>
</svelte:head>

<Navbar />

{@render children?.()}

<Footer />