<script>
    import { page } from '$app/stores';
    import {
        Content,
        Header,
        HeaderNav,
        HeaderNavItem,
    } from "carbon-components-svelte";

</script>

<svelte:head>
    <link
            rel="stylesheet"
            href="https://unpkg.com/carbon-components-svelte/css/white.css"
    />
</svelte:head>


<Header platformName="isumaru" href="/">
    <HeaderNav>
        <HeaderNavItem href="/" text="TOP" isSelected={$page.url.pathname == "/"}/>
        <HeaderNavItem href="/mysql" text="MYSQL" isSelected={$page.url.pathname == "/mysql"}/>
        <HeaderNavItem href="/setting" text="SETTING" isSelected={$page.url.pathname == "/setting"}/>
    </HeaderNav>
</Header>

<Content>
    <slot />
</Content>
