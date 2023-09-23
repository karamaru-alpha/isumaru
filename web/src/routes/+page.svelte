<script lang="ts">
    import {onMount} from "svelte";
    import {
        Button,
        DataTable,
        TextInput,
        Link,
        Tooltip,
    } from "carbon-components-svelte";
    import Play from "carbon-icons-svelte/lib/Play.svelte";
    import Launch from "carbon-icons-svelte/lib/Launch.svelte";
    import Warning from "carbon-icons-svelte/lib/Warning.svelte";


    enum targetType {
        slowQueryLog = 1,
    }
    enum statusType {
        progress = 1,
        success = 2,
        failure = 3,
    }

    let entries: {
        id: string;
        targets: {
            id: string;
            type: targetType;
            statusType: statusType;
            errorMessage: string;
        }[];
    }[] = [];

    onMount(async () => {
        const res = await fetch("http://localhost:8000/group")
        const data = await res.json()
        entries = data.entries;
    })

    async function collect() {
        try {
            await fetch("http://localhost:8000/group/collect", {
                method: "POST",
            })
        } catch (e) {
            console.log(e)
        }
    }
</script>

<div class="justify-space-between">
    <TextInput readonly labelText="Webhook URL" value="localhost:8000/group/collect"/>
    <Button icon={Play} on:click={collect}>Collect</Button>
</div>

<br />


{#each entries as entry}
    <p>{new Date(entry.id * 1000).toLocaleString()}</p>
    <DataTable
        headers={[
            { key: "id", value: "ID", width: "33%" },
            { key: "type", value: "Type", width: "33%" },
            { key: "statusType", value: "Status", width: "33%" },
        ]}
        rows={entry.targets}
        size="short"
    >
        <svelte:fragment slot="cell" let:cell let:row>
            {#if cell.key == "id" && row.statusType == statusType.success}
                <Link icon={Launch} href="/mysql/{entry.id}/{cell.value}">
                    {cell.value}
                </Link>
            {:else if cell.key == "type"}
                {targetType[cell.value]}
            {:else if cell.key == "statusType"}
                {#if cell.value == statusType.failure}
                    <Tooltip triggerText={statusType[cell.value]} icon={Warning}>
                        {row.errorMessage}
                    </Tooltip>
                {:else}
                    {statusType[cell.value]}
                {/if}
            {:else}
                {cell.value}
            {/if}
        </svelte:fragment>
    </DataTable>
    <br />
{/each}

<style>
    .justify-space-between {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
</style>
