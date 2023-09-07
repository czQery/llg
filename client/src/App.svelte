<script lang="ts">
    import Graph from "./lib/Graph.svelte";
    import Details from "./lib/Details.svelte";

    /*const init = async () => {

    }

    (async () => {
        while (true) {
            if (document.getElementById("app-end")) {
                init().then();
                break;
            }
            await sleep(50);
        }
    })();*/
</script>

<main>
    <header id="header">
        <h1>Login/Logoff Graph</h1>
        <form>
            <div>
                <label for="month">Month:</label>
                <input type="month" id="month" name="month" min="2020-00" value="2023-01"/>
            </div>
            <div>
                <button>render</button>
            </div>
        </form>
    </header>
    <div id="wrapper">
        <Details/>
        <Graph/>
    </div>
    <footer id="footer">
        <div style="background-color: blue; width: 60%">30h</div>
        <div style="background-color: red; width: 40%">20h</div>
    </footer>
    <div id="app-end" style="display: block;opacity: 0;width: 100%;height: 0px;position: relative;"></div>
</main>

<style>
    #header {
        height: 100%;
        display: grid;
        grid-column: 1/3;
        grid-row: 1/2;
        grid-template-columns: 1fr 1fr;
        grid-template-rows: 1fr;
        overflow: hidden;
        margin-bottom: 20px;
    }

    #header h1 {
        color: #FFF;
        margin: 0;
    }

    #header form {
        display: flex;
        flex-direction: column;
        gap: 10px;
        color: #FFF;
        text-align: right;
    }

    #header form div {
        display: flex;
        justify-content: end;
        gap: 10px;
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
        width: 140px;
        background-color: #3b60dc;
        font-size: 18px;
        color: #FFF;
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
        display: flex;
        grid-column: 1/3;
        grid-row: 3/4;
        padding: 0 10px 10px 10px;
        background-color: #FFF;
        border-radius: 0 0 var(--rad) var(--rad);
        box-shadow: var(--shadow);
    }

    #footer div {
        height: 20px;
        display: inline-block;
        text-align: center;
        color: white;
        /*border-radius: var(--rad);*/
    }
</style>
