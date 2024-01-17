<script lang="ts">
    import {getDate, sleep} from "./lib/ts/helper";
    import Graph from "./lib/components/Graph.svelte";
    import {type dataSum, type infoSum, loadData, loadInfo} from "./lib/ts/api";
    import SumGraph from "./lib/components/SumGraph.svelte";
    import {dataStore, infoStore, type itemInput, itemInputStore} from "./lib/ts/global.js";
    import ItemInput from "./lib/components/ItemInput.svelte";
    import {parseSelection} from "./lib/ts/param";
    import {inputToString} from "./lib/ts/helper.js";

    let info: infoSum = {build: "", selected_items: 0, users: [], devices: []};
    let infoItemInput: itemInput[] = [];

    let aspElement: HTMLSpanElement;
    let dateInputElement: HTMLInputElement;
    let itemInputList: itemInput[] = [];

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

    const init = async () => {
        let urlParams = new URLSearchParams(window.location.search);
        let urlDate = urlParams.get("date");
        if (urlDate) {
            let dateP = Date.parse(urlDate);
            if (dateP) {
                dateInputElement.value = getDate(new Date(dateP));
            }
        } else {
            dateInputElement.value = getDate(new Date());
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
        let dateP = Date.parse(dateInputElement.value);
        if (!dateP) {
            dateInputElement.value = getDate(new Date());
            return;
        }
        let date = new Date(dateP);
        dateInputElement.value = getDate(date);

        const url = new URL(window.location.toString());
        url.searchParams.set("date", getDate(date));

        url.searchParams.set("users", inputToString(itemInputList, "user"));
        url.searchParams.set("devices", inputToString(itemInputList, "device"));

        window.history.replaceState(null, "", url.toString());

        if (itemInputList.length === 0) {
            dataStore.set({} as dataSum);
            return;
        }

        let param = "?date=" + encodeURIComponent(dateInputElement.value)
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
                       on:change={render}
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
            <button on:click={() => {navigator.clipboard.writeText(aspElement.innerText)}}>copy
                admin path
            </button>
        </div>
    </header>
    <div id="wrapper">
        <!--<Details/>-->
        <Graph aspElement={aspElement} itemInputList={itemInputList}/>
    </div>
    <SumGraph/>
    <footer id="footer">
        <span>Created by <a href="https://qery.cz/l/g_llg">Štěpán Aubrecht</a></span>
        <span>build: <a href="https://github.com/czQery/llg/releases">{info.build}</a></span>
    </footer>
    <div id="app-end" style="display: block;opacity: 0;width: 100%;height: 0px;position: relative;"></div>
</main>

<style>
    :global(.sv-control) {
        border: 1px solid #FFF !important;
        border-radius: var(--rad) !important;
        min-height: 30px !important;
        height: 30px !important;
        color: #000 !important;
        outline: none !important;
    }

    :global(.sv-control:hover) {
        cursor: text;
    }

    :global(.sv-control.is-active) {
        border: 1px solid var(--cb) !important;
        outline: none !important;
    }

    :global(.indicator-icon) {
        transition: var(--transition) color linear !important;
    }

    :global(.indicator-icon:hover) {
        cursor: pointer !important;
    }

    :global(.sv-content.sv-input-row) {
        transition: var(--transition) font-size linear;
    }

    :global(.sv-content.sv-input-row:has(.inputBox:focus)) {
        font-size: 0 !important;
    }

    :global(.inputBox) {
        font-size: 16px !important;
        color: #000 !important;
    }

    :global(.inputBox:focus::placeholder) {
        color: #000 !important;
        opacity: 0 !important;
    }

    :global(.inputBox::placeholder) {
        color: #000 !important;
        opacity: 0.4 !important;
    }

    :global(.sv-item) {
        margin: 2px !important;
        border-radius: var(--rad) !important;
    }

    :global(.sv-item-btn) {
        border-radius: var(--rad) !important;
        transition: var(--transition) all linear !important;
    }

    :global(.sv-item-btn:hover) {
        background-color: var(--sv-item-btn-bg, var(--sv-item-selected-bg)) !important;
        filter: var(--hover);
    }

    :global(.sv-item-btn:active) {
        animation: reverse click-scale 100ms !important;
    }

    :global(.sv-item-btn svg) {
        fill: var(--c0) !important;
    }

    :global(.sv-dropdown) {
        border-radius: var(--rad) !important;
    }

    :global(.sv-dropdown .has-multiSelection) {
        flex-direction: column;
    }

    :global(.sv-dropdown-scroll.is-empty) {
        height: 4px;
    }

    :global(.sv-dropdown-scroll .sv-group-header) {
        color: #000;
        text-align: left !important;
    }

    :global(.sv-control.is-active > .has-multiSelection > div) {
        width: 100%;
    }

    :global(.sv-item-content) {
        border-radius: var(--rad) !important;
    }

    :global(.sv-item-content div) {
        display: flex;
        width: 100%;
        max-width: 183px;
        overflow: hidden;
    }

    :global(.sv-dd-item .sv-item) {
        transition: var(--transition) all linear !important;
    }

    :global(.sv-dd-item .sv-item-content div) {
        max-width: 192px !important;
    }

    :global(.sv-item-content div svg) {
        height: 19px;
        width: 19px;
        margin: 0 6px 0 0;
        padding: 2px 0;
        border-radius: var(--rad);
        background-color: var(--bg);
        stroke: white;
        display: inline-block;
        min-width: 19px;
    }

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
        opacity: 0.2;
        transition: var(--transition) opacity linear;
    }

    #header form div input::-webkit-calendar-picker-indicator:hover {
        opacity: 0.4;
        cursor: pointer;
    }

    #header form div input:focus {
        border: 1px solid var(--cb);
    }

    :global(.svelecte-control.svelecte) {
        height: 30px;
        width: 230px !important;
        flex-grow: 0 !important;
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
        display: flex;
        justify-content: space-between;
        color: white;
        font-size: 15px;
        text-align: right;
    }

    @media screen and (max-width: 809px) {
        #header {
            grid-template-columns: 1fr;
            grid-template-rows: min-content repeat(2, 70px);
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
            text-align: left;
            margin: 0;
        }

        #header-asp button {
            margin: 0;
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
