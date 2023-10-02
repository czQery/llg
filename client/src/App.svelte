<script lang="ts">
    import {getDate, getPaletteColor, sleep} from "./lib/ts/helper";
    import Graph from "./lib/components/Graph.svelte";
    import {type infoSum, type dataSum, loadData, loadInfo} from "./lib/ts/api";
    import SumGraph from "./lib/components/SumGraph.svelte";
    import {infoStore, type userInput, userInputStore, dataStore} from "./lib/ts/global.js";
    import UserInput from "./lib/components/UserInput.svelte";

    let info: infoSum = {build: "", users: [], selected_users: 0};
    let infoUserInput: userInput[] = [];

    let aspElement: HTMLSpanElement;
    let dateInputElement: HTMLInputElement;
    let userInputList: userInput[] = [];

    infoStore.subscribe((value: infoSum) => {
        if (value) {
            info = value;
            infoUserInput = info.users.map((u: string) => ({text: u, color: ""}))
        }
    });

    userInputStore.subscribe((value: userInput[]) => {
        if (value) {
            userInputList = value;
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

        let urlUsers = urlParams.get("users");
        let urlUsersInput: userInput[] = [];
        if (urlUsers) {
            let urlUsersList: string[] = urlUsers.split(",").filter((u: string) => {
                for (const uInfo of info.users) {
                    if (uInfo == u) {
                        return true;
                    }
                }

                return false;
            })

            urlUsersInput = urlUsersList.map((u: string) => {
                return {text: u, color: ""};
            });
        } else {
            urlUsersInput = info.users.slice(0, info.selected_users).map((u: string) => {
                return {text: u, color: ""}
            });
        }

        for (const u of urlUsersInput) {
            let color = getPaletteColor();
            let ttl = 0;

            for (let i = 0; i < urlUsersInput.length; i++) {
                if (urlUsersInput[i].color == color && urlUsersInput[i].text != u.text) {
                    color = getPaletteColor();

                    if (ttl > urlUsersInput.length * urlUsersInput.length) {
                        continue;
                    }

                    ttl = ttl + 1;
                    i = -1;
                }
            }
            u.color = color;
        }

        userInputStore.set(urlUsersInput);

        await render();

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
        url.searchParams.set("users", (userInputList.map((u: userInput) => {
            return u.text
        })).toString());
        window.history.pushState(null, "", url.toString());

        if (userInputList.length === 0) {
            dataStore.set({} as dataSum);
            return;
        }

        let param = "?date=" + encodeURIComponent(dateInputElement.value) + "&users=" + encodeURIComponent((userInputList.map((u: userInput) => {
            return u.text
        })).toString())

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
                <label for="month">Month:</label>
                <input type="month" id="month" name="month" min="2000-00"
                       on:change={render}
                       bind:this={dateInputElement}/>
            </div>
            <div>
                <label for="users">Users:</label>
                <UserInput userList={infoUserInput}/>
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
        <Graph aspElement={aspElement} userInputList={userInputList}/>
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

    :global(.sv-item) {
        margin: 2px !important;
    }

    :global(.sv-dropdown .has-multiSelection) {
        flex-direction: column;
    }

    :global(.sv-dropdown-scroll.is-empty) {
        height: 4px;
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

    #header a {
        grid-column: 1/2;
        grid-row: 1/2;
        color: #FFF;
        font-weight: 600;
        text-align: left;
        height: 70px;
        line-height: 70px;
        margin: 0;
        font-size: 80px;
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
