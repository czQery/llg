<script lang="ts">
    import {getDate, sleep} from "./lib/ts/helper";
    import Graph from "./lib/components/Graph.svelte";
    import {type infoSum, loadData, loadInfo} from "./lib/ts/api";
    import SumGraph from "./lib/components/SumGraph.svelte";
    import {infoStore} from "./lib/ts/global.js";
    import {MultiSelect} from "svelte-multiselect";

    let info: infoSum = {build: "", users: []};

    infoStore.subscribe(async (value: infoSum) => {
        if (value) info = value;
    });

    let aspTextElement: HTMLSpanElement;
    let formMonthElement: HTMLInputElement;

    const init = async () => {
        let urlParams = new URLSearchParams(window.location.search);
        let urlDate = urlParams.get("date");
        if (urlDate) {
            let dateP = Date.parse(urlDate);
            if (dateP) {
                formMonthElement.value = getDate(new Date(dateP));
            }
        } else {
            formMonthElement.value = getDate(new Date());
        }

        await loadData("?date=" + formMonthElement.value);

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

        formMonthElement.value = getDate(date);

        const url = new URL(window.location.toString());
        url.searchParams.set("date", getDate(date));
        window.history.pushState(null, "", url.toString());

        loadData("?date=" + formMonthElement.value).then();
    }

    (async () => {
        while (true) {
            if (document.getElementById("app-end")) {
                let ok = await loadInfo();

                if (!ok) {
                    await sleep(3000);
                    continue
                }

                init().then();
                break;
            }
            await sleep(50);
        }
    })();
</script>

<main>
    <header id="header">
        <h1>LLG</h1>
        <form>
            <div>
                <label for="month">Month:</label>
                <input type="month" id="month" name="month" min="2000-00" on:change={render}
                       bind:this={formMonthElement}/>
            </div>
            <div>
                <label for="users">Users:</label>
                <MultiSelect id="users" on:change={render}
                             --sms-options-bg="#FFF"
                             --sms-border="2px solid #FFF"
                             --sms-border-radius="5px"
                             --sms-focus-border="2px solid #3b60dc"
                             --sms-bg="#FFF"
                             --sms-text-color="#000"
                             allowUserOptions="append"
                             selected={info.users.slice(0, 3)}
                             options={info.users}/>
            </div>
        </form>
        <div id="header-asp">
            <span bind:this={aspTextElement}></span>
            <button on:click={() => {navigator.clipboard.writeText(aspTextElement.innerText)}}>copy
                admin path
            </button>
        </div>
    </header>
    <div id="wrapper">
        <!--<Details/>-->
        <Graph aspElement={aspTextElement}/>
    </div>
    <SumGraph/>
    <footer id="footer">
        <span>Created by <a href="https://qery.cz/l/g_llg">Štěpán Aubrecht</a></span>
        <span>build: <a href="https://github.com/czQery/llg/releases">{info.build}</a></span>
    </footer>
    <div id="app-end" style="display: block;opacity: 0;width: 100%;height: 0px;position: relative;"></div>
</main>

<style>
    #header {
        height: 100%;
        display: grid;
        grid-column: 1/3;
        grid-row: 1/2;
        grid-template-columns: min-content 1fr max-content;
        grid-template-rows: 70px;
        gap: 10px;
        padding: 0 0 20px 0;
    }

    #header h1 {
        grid-column: 1/2;
        grid-row: 1/2;
        color: #FFF;
        font-weight: 600;
        text-align: center;
        height: 70px;
        line-height: 70px;
        margin: 0;
        font-size: 80px;
        border-right: 3px solid white;
        padding-right: 7px;
    }

    #header form {
        grid-column: 2/3;
        grid-row: 1/2;
        display: flex;
        flex-direction: column;
        gap: 10px;
        color: #FFF;
        text-align: right;
    }

    #header form div {
        display: flex;
        gap: 10px;
        flex-direction: row;
    }

    #header form div label {
        font-size: 25px;
        height: 30px;
        width: 75px;
        line-height: 30px;
        text-align: left;
    }

    #header form div input {
        height: 30px;
        width: 180px;
        border: 2px solid #FFF;
        font-size: 18px;
        border-radius: var(--rad);
        text-align: center;
        outline: none;
    }

    #header form div input:focus {
        border: 2px solid #3b60dc;
    }

    #header-asp {
        display: flex;
        flex-direction: column;
        gap: 10px;
        grid-column: 3/4;
        grid-row: 1/2;
    }

    #header-asp span {
        color: white;
        text-align: right;
        font-size: 25px;
        margin: auto 0 0 auto;
        overflow: hidden;
    }

    #header-asp button {
        height: 30px;
        max-width: 176.5px;
        width: 100%;
        background-color: var(--cb);
        font-size: 18px;
        color: #FFF;
        margin: 0 0 0 auto;
        white-space: nowrap;
        padding: 0 10px;
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

    #footer {
        padding: 10px 0 0 0;
        width: 100%;
        display: flex;
        justify-content: space-between;
        color: white;
        font-size: 15px;
        text-align: right;
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

        #header-asp span {
            grid-column: 1/2;
            grid-row: 3/4;
            text-align: left;
            margin: 0;
        }

        #header-asp button {
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
