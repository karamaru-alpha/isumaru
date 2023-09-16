<script lang="ts">
    import { DataTable } from "carbon-components-svelte";
    import { page } from '$app/stores';
    import {onMount} from "svelte";


    const id = $page.params.id;

    let headers: {
        key: string,
        value: string
    }[];
    let rows: {}[];

    onMount(async () => {
        try {
            const res = await fetch("http://localhost:8000/mysql/" + id)
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


<h1>Mysql {id}</h1>
<br />

<DataTable
    sortable
    headers={headers}
    rows={rows}
/>
