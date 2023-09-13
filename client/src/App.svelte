<script lang="ts">
    import {getDate, sleep} from "./lib/ts/helper";
    import Details from "./lib/components/Details.svelte";
    import Graph from "./lib/components/Graph.svelte";
    import {loadData} from "./lib/ts/api";
    import SumGraph from "./lib/components/SumGraph.svelte";

    let aspTextElement: HTMLSpanElement;
    let formMonthElement: HTMLInputElement;

    const init = async () => {
        let urlParams = new URLSearchParams(window.location.search);
        let urlDate = urlParams.get("date");
        if (urlDate) {
            let date = Date.parse(urlDate);
            if (date) {
                formMonthElement.value = urlDate;
            }
        } else {
            formMonthElement.value = getDate(new Date());
        }

        loadData("?date="+formMonthElement.value).then();

        setTimeout(() => {
            document.getElementById("loading")!.style.display = "none";
            document.getElementById("app")!.style.display = "block";
        }, 50)
    }

    const render = async () => {
        let dateP = Date.parse(formMonthElement.value);
        if (!dateP) {
            formMonthElement.value = getDate(new Date());
            return;
        }

        let date = new Date(dateP)

        const url = new URL(window.location.toString());
        url.searchParams.set("date", getDate(date));
        window.history.pushState(null, "", url.toString());

        loadData("?date="+formMonthElement.value).then();
    }

    (async () => {
        while (true) {
            if (document.getElementById("app-end")) {
                init().then();
                break;
            }
            await sleep(50);
        }
    })();
</script>

<main>
    <header id="header">
        <h1>Login/Logoff Graph</h1>
        <form>
            <label for="month">Month:</label>
            <input type="month" id="month" name="month" min="2000-00" bind:this={formMonthElement}/>
            <button on:click|preventDefault={render}>render</button>
        </form>
        <span id="header-asp-text" bind:this={aspTextElement}></span>
        <button id="header-asp-btn" on:click={() => {navigator.clipboard.writeText(aspTextElement.innerText)}}>copy
            admin path
        </button>
    </header>
    <div id="wrapper">
        <Details/>
        <Graph aspElement={aspTextElement}/>
    </div>
    <SumGraph/>
    <div id="app-end" style="display: block;opacity: 0;width: 100%;height: 0px;position: relative;"></div>
</main>

<style>
    #header {
        height: 100%;
        display: grid;
        grid-column: 1/3;
        grid-row: 1/2;
        grid-template-columns: 1fr 1fr;
        grid-template-rows: 1fr 30px;
        overflow: hidden;
        gap: 10px;
        padding: 0 0 20px 0;
    }

    #header h1 {
        grid-column: 1/2;
        grid-row: 1/2;
        color: #FFF;
        font-size: 35px;
        font-weight: 600;
        margin: 0;
    }

    #header form {
        grid-column: 1/2;
        grid-row: 2/3;
        display: flex;
        gap: 10px;
        color: #FFF;
        text-align: right;
    }

    #header form label {
        font-size: 25px;
        height: 30px;
        line-height: 30px;
    }

    #header form input {
        height: 30px;
        width: 140px;
        border: none;
        font-size: 18px;
        border-radius: var(--rad);
        text-align: center;
    }

    #header form button {
        height: 30px;
        max-width: 140px;
        width: 100%;
        background-color: var(--cb);
        font-size: 18px;
        color: #FFF;
    }

    #header-asp-text {
        color: white;
        text-align: right;
        grid-column: 2/3;
        grid-row: 1/2;
        font-size: 25px;
        margin: auto 0 0 auto;
        overflow: hidden;
    }

    #header-asp-btn {
        grid-column: 2/3;
        grid-row: 2/3;
        height: 30px;
        max-width: 190px;
        width: 100%;
        background-color: var(--cb);
        font-size: 18px;
        color: #FFF;
        margin: 0 0 0 auto;
        white-space: nowrap;
    }

    #wrapper {
        display: block;
        height: 100%;
        grid-column: 1/3;
        grid-row: 2/3;
        padding: 10px;
        background-color: #FFF;
        border-radius: var(--rad) var(--rad) 0 0;
        box-shadow: var(--shadow);
    }

    @media screen and (max-width: 809px) {
        #header {
            grid-template-columns: 1fr max-content;
            grid-template-rows: 1fr repeat(2, 30px);
        }

        #header h1 {
            grid-column: 1/3;
            grid-row: 1/2;
        }

        #header form {
            grid-column: 1/3;
            grid-row: 2/3;
            display: flex;
            gap: 10px;
            color: #FFF;
            text-align: right;
        }

        #header-asp-text {
            grid-column: 1/2;
            grid-row: 3/4;
            text-align: left;
            margin: 0;
        }

        #header-asp-btn {
            grid-column: 2/3;
            grid-row: 3/4;
            margin: 0;
        }
    }

    @media screen and (max-width: 519px) {
        #header {
            padding: 0 0 10px 0;
        }
    }
</style>
