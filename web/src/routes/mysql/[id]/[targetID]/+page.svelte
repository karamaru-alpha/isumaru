<script lang="ts">
    import { page } from '$app/stores';
    import {
        DataTable,
        Button,
    } from "carbon-components-svelte";
    import Cube from "carbon-icons-svelte/lib/Cube.svelte";
    import {onMount} from "svelte";


    const id = $page.params.id;
    const currentTargetID = $page.params.targetID;
    let targetIDs: string[] = [];
    let headers: {
        key: string,
        value: string
    }[];
    let rows: {}[];

    onMount(async () => {
        try {
            const res = await fetch(`http://localhost:8000/mysql/${id}`)
            const json = await res.json();
            targetIDs = json.targetIDs;
        } catch (e) {
            console.log(e)
        }

        try {
            const res = await fetch(`http://localhost:8000/mysql/${id}/${currentTargetID}`)
            const tsv = await res.text();
            let lines = tsv.split('\n');
            if (tsv.endsWith('\n')) {
                lines = lines.slice(0, -1);
            }
            headers = lines[0].split('\t').map((v, i) => ({ key: i.toString(), value: v }));
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
    <Button kind={targetID == currentTargetID ? "primary" : "tertiary"} size="small" icon={Cube} on:click={window.location.href = `/mysql/${id}/${targetID}`}>{targetID}</Button>
{/each}
<br />
<br />

<p>Mysql ({new Date(id * 1000).toLocaleString()})</p>

<DataTable
    sortable
    headers={headers}
    rows={rows}
    size="short"
/>
