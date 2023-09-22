<script lang="ts">
    import {onMount} from "svelte";
    import {
        Button,
        DataTable,
        Link,
    } from "carbon-components-svelte";
    import Launch from "carbon-icons-svelte/lib/Launch.svelte";
    import Play from "carbon-icons-svelte/lib/Play.svelte";

    let entries: {
        id: string;
        time: string;
        targets: string[];
    }[];
    onMount(async () => {
        try {
            const res = await fetch("http://localhost:8000/mysql")
            const data = await res.json()
            entries = data.entries.map(e => {
                return {
                    id: e.id,
                    time: new Date(e.unixTime * 1000).toLocaleString(),
                    targets: e.targetIDs,
                }
            });
        } catch (e) {
            console.log(e)
        }
    });

    function collect() {
        fetch("http://localhost:8000/mysql/collect", {
            method: "POST",
        });
    }

</script>

<div class="justify-flex-end">
    <Button kind="tertiary" icon={Play} size="small" on:click={collect}>Collect</Button>
</div>

<DataTable
    sortable
    headers={[
        { key: "id", value: "ID" },
        { key: "targets", value: "Targets" },
        { key: "time", value: "Time" },
    ]}
    rows={entries}
    size="short"
>
    <svelte:fragment slot="cell" let:cell let:row>
        {#if cell.key === "id"}
            <Link icon={Launch} href="/mysql/{cell.value}/{row.targets[0]}">
                {cell.value}
            </Link>
        {:else}
            {cell.value}
        {/if}
    </svelte:fragment>
</DataTable>

<style>
    .justify-flex-end {
        display: flex;
        align-items: center;
        justify-content: end;
    }
</style>
