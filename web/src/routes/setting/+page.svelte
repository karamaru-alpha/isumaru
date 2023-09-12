<script>
    import {onMount} from 'svelte';
    import {
        Button,
        Form,
        TextInput,
    } from "carbon-components-svelte";

    let seconds;
    let mainServerAddress;
    let accessLogPath;
    let mysqlServerAddress;
    let slowQueryLogPath;

    onMount(async () => {
        try {
            const res = await fetch("http://localhost:8000/setting")
            const data = await res.json()
            seconds = data.seconds
            mainServerAddress = data.mainServerAddress
            accessLogPath = data.accessLogPath
            mysqlServerAddress = data.mysqlServerAddress
            slowQueryLogPath = data.slowQueryLogPath
        } catch (e) {
            console.log(e)
        }
    });

    async function onSubmit(e) {
        e.preventDefault();
        try {
            await fetch("http://localhost:8000/setting", {
                method: "POST",
                body: JSON.stringify({
                    seconds,
                    mainServerAddress,
                    accessLogPath,
                    mysqlServerAddress,
                    slowQueryLogPath,
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

<Form on:submit={onSubmit}>

    <TextInput type="number" labelText="メトリクス時間(s)" placeholder="75" bind:value={seconds} />
    <br />

    <TextInput type="string" labelText="メインサーバーIPアドレス" placeholder="127.0.0.1" bind:value={mainServerAddress} />
    <TextInput type="string" labelText="アクセスログのPath" placeholder="/var/log/nginx/access.log" bind:value={accessLogPath} />
    <br />

    <TextInput type="string" labelText="MysqlサーバーIPアドレス" placeholder="127.0.0.1" bind:value={mysqlServerAddress} />
    <TextInput type="string" labelText="スロークエリログのPath" placeholder="/var/log/mysql/slow-query.log" bind:value={slowQueryLogPath} />
    <br />

    <Button type="submit">Submit</Button>
</Form>
