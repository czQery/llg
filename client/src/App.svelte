<script lang="ts">
    import {getDate, getPaletteColor, sleep} from "./lib/ts/helper";
    import Graph from "./lib/components/Graph.svelte";
    import {type infoSum, loadData, loadInfo} from "./lib/ts/api";
    import SumGraph from "./lib/components/SumGraph.svelte";
    import {type activeUser, activeUserStore, infoStore} from "./lib/ts/global.js";

    import Svelecte, {addFormatter} from "svelecte";

    let info: infoSum = {build: "", users: []};
    let infoInput: { text: string }[] = [];
    let users: activeUser[] = [];

    let aspTextElement: HTMLSpanElement;
    let formMonthElement: HTMLInputElement;
    let formUsersList: string[];

    infoStore.subscribe(async (value: infoSum) => {
        if (value) {
            info = value;
            infoInput = info.users.map((u: string) => ({"text": u}))
        }
    });
    activeUserStore.subscribe(async (value: activeUser[]) => {
        if (value) users = value;
    });

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

        let urlUsers = urlParams.get("users");
        if (urlUsers) {
            formUsersList = urlUsers.split(",");
        } else {
            formUsersList = info.users.slice(0, info.selected_users);
        }

        await render();

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
        let date = new Date(dateP);
        formMonthElement.value = getDate(date);

        const url = new URL(window.location.toString());
        url.searchParams.set("date", getDate(date));
        url.searchParams.set("users", formUsersList.toString());
        window.history.pushState(null, "", url.toString());

        await loadData("?date=" + encodeURIComponent(formMonthElement.value) + "&users=" + encodeURIComponent(formUsersList.toString()));
    }

    const userColor = (item: { text: string }, isSelected: boolean) => {
        if (isSelected) {
            let color = getPaletteColor();
            let ttl = 0;

            for (let j = 0; j < users.length; j++) {
                if (users[j].color == color) {
                    color = getPaletteColor();

                    if (ttl > users.length*4) {
                        continue;
                    }

                    ttl = ttl+1;
                    j = -1;
                }
            }

            let aUsers = users;
            aUsers.push({name: item.text, color: color, sum: -1})
            activeUserStore.set(aUsers);
            return '<div style="border-left: 4px solid ' + color + '; padding: 0 0 0 4px">' + item.text + '</div>';
        }

        return '<div>' + item.text + '</div>';
    }
    addFormatter("user-color", userColor);

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
                <input type="month" id="month" name="month" min="2000-00"
                       on:change={render}
                       bind:this={formMonthElement}/>
            </div>
            <div>
                <label for="users">Users:</label>
                <Svelecte inputId="users"
                          renderer="user-color"
                          options={infoInput}
                          multiple={true}
                          collapseSelection={true}
                          alwaysCollapsed={true}
                          searchable={true}
                          valueAsObject={false}
                          bind:value={formUsersList}
                          on:change={render}/>
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
    :global(.sv-control) {
        border: 1px solid #FFF !important;
        border-radius: var(--rad) !important;
        min-height: 30px !important;
        height: 30px !important;
        color: #000 !important;
        outline: none !important;
    }

    :global(.sv-control.is-active) {
        border: 1px solid var(--cb) !important;
        outline: none !important;
    }

    :global(.sv-content .inputBox) {
        color: #000;
        width: 1px;
        height: 0 !important;
    }

    :global(.sv-dropdown .has-multiSelection) {
        flex-direction: column;
    }

    :global(.sv-control.is-active > .has-multiSelection > div) {
        width: 100%;
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
        text-align: left;
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
        border: 1px solid #FFF;
        font-size: 16px;
        border-radius: var(--rad);
        text-align: left;
        outline: none;
        padding: 0 0 0 5px;
    }

    #header form div input:focus {
        border: 1px solid var(--cb);
    }

    :global(.svelecte-control.svelecte) {
        height: 30px;
        width: 180px !important;
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

        #header h1 {
            text-align: left;
            border-right: none;
        }

        #header form {
            grid-column: 1/2;
            grid-row: 2/3;
            display: flex;
            gap: 10px;
            color: #FFF;
            text-align: right;
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
