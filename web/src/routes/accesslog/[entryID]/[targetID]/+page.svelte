<script lang="ts">
    import { page } from '$app/stores';
    import {
        DataTable,
        Button,
    } from "carbon-components-svelte";
    import Cube from "carbon-icons-svelte/lib/Cube.svelte";
    import {onMount} from "svelte";

    const entryID = $page.params.entryID;
    const currentTargetID = $page.params.targetID;
    let targetIDs: string[] = [];
    let headers: {
        key: string,
        value: string
    }[];
    let rows: {}[];

    onMount(async () => {
        try {
            const res = await fetch(`/api/accesslog/${entryID}/${currentTargetID}`)
            const json = await res.json();
            targetIDs = json.targetIDs;
            const tsv = atob(json.data);
            let lines = tsv.split('\n');
            if (tsv.endsWith('\n')) {
                lines = lines.slice(0, -1);
            }
            headers = lines[0].split('\t').map((v, i) => {
                if (i == 1) {
                    return {key: i.toString(), value: v }
                }
                return {key: i.toString(), value: v }
            });
            rows = lines.slice(1).map((line, index) => {
                let row: {} = {
                    id: index,
                };
                line.split('\t').forEach((v, i) => {
                    row[i.toString()] = v;
                });
                return row;
            });
        } catch (e) {
            console.log(e)
        }
    });
</script>

{#each targetIDs as targetID, index}
    <Button kind={targetID == currentTargetID ? "primary" : "tertiary"} size="small" icon={Cube} on:click={window.location.href = `/accesslog/${entryID}/${targetID}`}>{targetID}</Button>
{/each}
<br />
<br />

<p>AccessLog ({new Date(entryID * 1000).toLocaleString()})</p>

<DataTable
        sortable
        headers={headers}
        rows={rows}
        size="short"
/>
