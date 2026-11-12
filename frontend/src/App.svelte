<script lang="ts">
    import {getDate, sleep} from "./lib/ts/helper";
    import Graph from "./lib/components/Graph.svelte";
    import {type dataSum, type infoSum, loadData, loadInfo} from "./lib/ts/api";
    import SumGraph from "./lib/components/SumGraph.svelte";
    import {dataStore, infoStore, type itemInput, itemInputStore} from "./lib/ts/global.js";
    import ItemInput from "./lib/components/ItemInput.svelte";
    import {parseSelection} from "./lib/ts/param";
    import {inputToString} from "./lib/ts/helper.js";

    let info: infoSum = $state({build: "", selected_items: 0, users: [], devices: []});
    let infoItemInput: itemInput[] = $state([]);

    let data: dataSum = $state({dates: [], items: [], time: 0, lines: 0});

    let aspElement: HTMLSpanElement | undefined = $state();
    let dateInputElement: HTMLInputElement | undefined = $state();
    let itemInputList: itemInput[] = $state([]);

    infoStore.subscribe((value: infoSum) => {
        if (value) {
            info = value;
            let list: itemInput[] = [];

            if (info.users) {
                list.push(...info.users.map((u: string) => (<itemInput>{type: "user", value: u, color: ""})));
            }

            if (info.devices) {
                list.push(...info.devices.map((u: string) => (<itemInput>{type: "device", value: u, color: ""})));
            }

            infoItemInput = list;
        }
    });

    itemInputStore.subscribe((value: itemInput[]) => {
        if (value) {
            itemInputList = value;
            render();
        }
    });

    dataStore.subscribe((value: dataSum) => {
        if (value) {
            data = value;
        }
    });

    const init = async () => {
        let urlParams = new URLSearchParams(window.location.search);
        let urlDate = urlParams.get("date");
        if (urlDate) {
            let dateP = Date.parse(urlDate);
            if (dateP) {
                dateInputElement!.value = getDate(new Date(dateP));
            }
        } else {
            dateInputElement!.value = getDate(new Date());
        }

        let list: itemInput[] = [];

        list.push(...parseSelection("user", urlParams, info))
        list.push(...parseSelection("device", urlParams, info))

        itemInputStore.set(list);

        setTimeout(() => {
            document.getElementById("loading")!.style.display = "none";
            document.getElementById("app")!.style.display = "block";
        }, 50)
    }

    let renderLast: string = "";
    const render = async () => {
        let dateP = Date.parse(dateInputElement!.value);
        if (!dateP) {
            dateInputElement!.value = getDate(new Date());
            return;
        }
        let date = new Date(dateP);
        dateInputElement!.value = getDate(date);

        const url = new URL(window.location.toString());
        url.searchParams.set("date", getDate(date));

        url.searchParams.set("users", inputToString(itemInputList, "user"));
        url.searchParams.set("devices", inputToString(itemInputList, "device"));

        window.history.replaceState(null, "", url.toString());

        if (itemInputList.length === 0) {
            dataStore.set({dates: [], items: [], time: 0, lines: 0} as dataSum);
            return;
        }

        let param = "?date=" + encodeURIComponent(dateInputElement!.value)
            + "&users="
            + encodeURIComponent(inputToString(itemInputList, "user"))
            + "&devices="
            + encodeURIComponent(inputToString(itemInputList, "device"));

        if (renderLast === param) {
            return;
        }

        renderLast = param

        await loadData(param);
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
        <a href="/">LLG</a>
        <form>
            <div>
                <label for="date">Month:</label>
                <input id="date" type="month" name="month" min="2000-00"
                       onchange={render}
                       bind:this={dateInputElement}/>
            </div>
            <div>
                <label for="items" style="letter-spacing: 1.5px">Items:</label>
                <ItemInput op={[
                    {
                        groupHeader: "Users",
                        items: infoItemInput.filter((i) => i.type === "user")
                    },
                    {
                        groupHeader: "Devices",
                        items: infoItemInput.filter((i) => i.type === "device")
                    }
                ]}/>
            </div>
        </form>
        <div id="header-asp">
            <span bind:this={aspElement}></span>
            <button onclick={() => {navigator.clipboard.writeText(aspElement!.innerText)}}>copy
                admin path
            </button>
        </div>
    </header>
    <div id="wrapper">
        <Graph aspElement={aspElement} itemInputList={itemInputList}/>
    </div>
    <SumGraph/>
    <footer id="footer">
        <span>Created by <a href="https://qery.cz">Štěpán Aubrecht</a></span>
        <div id="stats">
            <span>{data.time / 1_000_000 + " ms"}</span>
            <span>|</span>
            <span>{data.lines + " lines"}</span>
        </div>
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

    #header a {
        grid-column: 1/2;
        grid-row: 1/2;
        color: #FFF;
        font-weight: 600;
        text-align: left;
        height: 70px;
        line-height: 70px;
        margin: 0;
        max-width: max-content;
        font-size: 80px;
        overflow: hidden;
    }

    #header form {
        grid-column: 2/3;
        grid-row: 1/2;
        display: flex;
        flex-direction: column;
        gap: 10px;
        color: #FFF;
        text-align: left;
        border-left: 3px solid white;
        padding-left: 13px;
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
        display: flex;
        height: 30px;
        width: 230px;
        border: 1px solid #FFF;
        font-size: 16px;
        border-radius: var(--rad);
        text-align: left;
        outline: none;
        padding: 0 8px 0 5px;  /* 8px for edge calendar icon */
    }

    #header form div input:hover {
        cursor: text;
    }

    #header form div input::-webkit-calendar-picker-indicator {
        opacity: 0.265;
        transition: var(--transition) opacity linear;
        margin: 0 0 0 55px;
    }

    #header form div input::-webkit-calendar-picker-indicator:hover {
        opacity: 0.81;
        cursor: pointer;
    }

    #header form div input:focus {
        border: 1px solid var(--cb);
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
        height: 30px;
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
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        color: white;
        font-size: 15px;
    }

    #footer span {
        width: max-content;
    }

    #footer span:first-child {
        margin: 0 auto 0 0;
    }

    #footer span:last-child {
        margin: 0 0 0 auto;
    }

    #footer #stats {
        display: flex;
        gap: 5px;
        opacity: 0.2;
        margin: 0 auto;
    }

    #footer #stats span {
        display: flex;
    }

    #footer #stats span:first-child {
        width: 100px;
        justify-content: end;
    }

    #footer #stats span:last-child {
        width: 100px;
        justify-content: start;
    }

    @media screen and (max-width: 809px) {
        #header {
            grid-template-columns: 1fr;
            grid-template-rows: min-content 70px 30px;
        }

        #header form {
            grid-column: 1/2;
            grid-row: 2/3;
            display: flex;
            gap: 10px;
            color: #FFF;
            text-align: right;
            border-left: none;
            padding-left: 5px;
        }

        #header-asp {
            grid-column: 1/2;
            grid-row: 3/4;
        }

        #header-asp span {
            margin: 0;
            display: none;
        }

        #header-asp button {
            margin: 0;
            max-width: none;
        }

        #footer #stats {
            width: 0;
            height: 0;
            opacity: 0;
        }
    }

    @media screen and (max-width: 519px) {
        #header {
            padding: 0 0 10px 0;
        }

        #footer {
            padding: 5px 0 0 0;
        }
    }
</style>
