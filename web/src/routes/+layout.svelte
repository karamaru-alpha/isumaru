<script>
    import { page } from '$app/stores';
    import {
        Content,
        Header,
        HeaderNav,
        HeaderNavItem,
    } from "carbon-components-svelte";
    import { SvelteToast } from '@zerodevx/svelte-toast'
</script>

<svelte:head>
    <link
        rel="stylesheet"
        href="https://unpkg.com/carbon-components-svelte/css/white.css"
    />
</svelte:head>

<SvelteToast />

<Header platformName="isumaru" href="/">
    <HeaderNav>
        <HeaderNavItem href="/" text="TOP" isSelected={$page.url.pathname == "/"}/>
        <HeaderNavItem href="/setting" text="SETTING" isSelected={$page.url.pathname == "/setting"}/>
    </HeaderNav>
</Header>

<Content>
    <slot />
</Content>
