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
    import {success, error} from "../../lib/toast";
    import {targetType} from '../../lib/enum'

    let targets: {
        id: string
        type: targetType,
        url: string,
        path: string,
        duration: number,
    }[] = [];
    let slpConfig: string;
    let alpConfig: string;

    onMount(async () => {
        try {
            const res = await fetch("/api/setting")
            if (!res.ok) {
                const err = await res.json();
                throw new Error(err.message);
            }

            const data = await res.json()
            targets = data.targets
            slpConfig = data.slpConfig
            alpConfig = data.alpConfig
        } catch (e) {
            error(`Failure: ${e.message}`);
        }
    });

    function addTarget() {
        targets = [
            ...targets,
            {
                id: "isu1",
                type: targetType.slowQueryLog,
                url: "http://localhost:19000",
                path: "/var/log/mysql/slow-query.log",
                duration: 70,
            }
        ]
    }

    async function saveTarget() {
        try {
            const res = await fetch("/api/setting/target", {
                method: "POST",
                body: JSON.stringify({
                    targets,
                }),
                headers: {
                    "Content-Type": "application/json"
                }
            })
            if (!res.ok) {
                const err = await res.json();
                throw new Error(err.message);
            }

            success("Succeeded");
        } catch (e) {
            error(`Failure: ${e.message}`);
        }
    }

    async function saveSlpConfig() {
        try {
            const res = await fetch("/api/setting/slp", {
                method: "POST",
                body: JSON.stringify({
                    slpConfig,
                }),
                headers: {
                    "Content-Type": "application/json"
                }
            })
            if (!res.ok) {
                const err = await res.json();
                throw new Error(err.message);
            }

            success("Succeeded");
        } catch (e) {
            error(`Failure: ${e.message}`);
        }
    }

    async function saveAlpConfig() {
        try {
            const res = await fetch("/api/setting/alp", {
                method: "POST",
                body: JSON.stringify({
                    alpConfig,
                }),
                headers: {
                    "Content-Type": "application/json"
                }
            })
            if (!res.ok) {
                const err = await res.json();
                throw new Error(err.message);
            }

            success("Succeeded");
        } catch (e) {
            error(`Failure: ${e.message}`);
        }
    }
</script>

<p>Target Config</p>
{#each targets as target, index}
    <div class="justify-space-between">
        <Select
            inline
            selected={target.type}
            on:change={(e) => {
                target.type = Number(e.target.value)
            }}
        >
            <SelectItem value={targetType.slowQueryLog} text="SlowQueryLog" />
            <SelectItem value={targetType.accessLog} text="AccessLog" />
            <SelectItem value={targetType.pprof} text="PProf" />
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

<p style="margin-top: 10px">Alp Config</p>
<TextArea
    rows={8}
    bind:value={alpConfig}
    label="Alp Config"
/>
<br />
<div class="justify-flex-end">
    <Button kind="tertiary" icon={WatsonHealthSubVolume} size="small" on:click={saveAlpConfig}>Save</Button>
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
