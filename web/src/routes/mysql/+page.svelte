<script lang="ts">
    import { DataTable, Link } from "carbon-components-svelte";
    import Launch from "carbon-icons-svelte/lib/Launch.svelte";
    import {onMount} from "svelte";

    let entries: {
        id: string;
        time: string;
    }[];
    onMount(async () => {
        try {
            const res = await fetch("http://localhost:8000/mysql")
            const data = await res.json()
            entries = data.entries.map(e => {
                return {
                    id: e.id,
                    time: new Date(e.unixTime * 1000).toLocaleString()
                }
            });
        } catch (e) {
            console.log(e)
        }
    });
</script>

<DataTable
    sortable
    headers={[
        { key: "id", value: "ID" },
        { key: "time", value: "Time" },
    ]}
    rows={entries}
>
    <svelte:fragment slot="cell" let:cell>
        {#if cell.key === "id"}
            <Link icon={Launch} href="/mysql/{cell.value}">
                {cell.value}
            </Link>
        {:else}
            {cell.value}
        {/if}
    </svelte:fragment>
</DataTable>
