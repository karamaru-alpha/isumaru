<script lang="ts">
    import {onMount} from 'svelte';
    import {
        Button,
        Select,
        SelectItem,
        TextArea,
        TextInput,
    } from "carbon-components-svelte";
    import Add from "carbon-icons-svelte/lib/Add.svelte";
    import WatsonHealthSubVolume from "carbon-icons-svelte/lib/WatsonHealthSubVolume.svelte";
    import TrashCan from "carbon-icons-svelte/lib/TrashCan.svelte";

    let targets: {
        id: string
        type: number,
        url: string,
        path: string,
        duration: number,
    }[] = [];
    let slpConfig: string;

    onMount(async () => {
        try {
            const res = await fetch("http://localhost:8000/setting")
            const data = await res.json()
            targets = data.targets
            slpConfig = data.slpConfig
        } catch (e) {
            console.log(e)
        }
    });

    function addTarget(e) {
        e.preventDefault();
        targets = [
            ...targets,
            {
                id: "isu1",
                type: 1,
                url: "http://localhost:8080",
                path: "/var/log/nginx/access.log",
                duration: 70,
            }
        ]
    }

    async function saveTarget(e) {
        e.preventDefault();
        try {
            await fetch("http://localhost:8000/setting/target", {
                method: "POST",
                body: JSON.stringify({
                    targets,
                }),
                headers: {
                    "Content-Type": "application/json"
                }
            })
        } catch (e) {
            console.log(e)
        }
    }
    async function saveSlpConfig(e) {
        e.preventDefault();
        try {
            await fetch("http://localhost:8000/setting/slp", {
                method: "POST",
                body: JSON.stringify({
                    slpConfig,
                }),
                headers: {
                    "Content-Type": "application/json"
                }
            })
        } catch (e) {
            console.log(e)
        }
    }
</script>

<p>Target Config</p>
{#each targets as target, index}
    <div class="justify-space-between">
        <Select
                class="flex-container"
                inline
                labelText="Type"
                selected={target.type.toString()}
                on:change={(e) => {
            target.type = Number(e.value)
        }}>
            <SelectItem value="1" text="AccessLog" />
            <SelectItem value="2" text="SlowQueryLog" />
        </Select>
        <Button kind="danger-tertiary" iconDescription="delete" icon={TrashCan} size="small" on:click={() => {
                targets = [...targets.slice(0, index), ...targets.slice(index + 1)];
        }} />
    </div>
    <TextInput
        size="sm"
        bind:value={target.id}
        inline
        labelText="ID"
        placeholder="isu1"
    />
    <TextInput
        size="sm"
        bind:value={target.url}
        inline
        labelText="URL"
        placeholder="http://localhost:8080"
    />
    <TextInput
        size="sm"
        bind:value={target.path}
        inline
        labelText="PATH"
        placeholder="/var/log/nginx/access.log"
    />
    <TextInput
        size="sm"
        bind:value={target.duration}
        inline
        type="number"
        labelText="Duration"
        placeholder="10"
    />
    <br />
{/each}

<div class="justify-flex-end">
    <Button kind="tertiary" iconDescription="add" icon={Add} size="small" on:click={addTarget} />
    <Button kind="tertiary" icon={WatsonHealthSubVolume} size="small" on:click={saveTarget}>Save</Button>
</div>

<p style="margin-top: 10px">Slp Config</p>
<TextArea
    rows={8}
    bind:value={slpConfig}
    label="Slp Config"
/>
<br />
<div class="justify-flex-end">
    <Button kind="tertiary" icon={WatsonHealthSubVolume} size="small" on:click={saveSlpConfig}>Save</Button>
</div>

<style>
    .justify-space-between {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
    .justify-flex-end {
        display: flex;
        align-items: center;
        justify-content: end;
    }
</style>
