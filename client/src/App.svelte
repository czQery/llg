<script lang="ts">
    import "chartjs-adapter-date-fns";
    import {Chart, LinearScale, TimeSeriesScale, BarController, BarElement, PointElement, Tooltip, Legend} from "chart.js";
    Chart.register(LinearScale, TimeSeriesScale, BarController, BarElement, PointElement, Tooltip, Legend);

    const formatTime = (value: number): string => {
        let h = Math.trunc(value / 60);
        let m = value % 60;

        let hs: string = h.toString();
        let ms: string = m.toString();

        if (h.toString().length === 1) {
            hs = "0" + h.toString();
        }

        if (m.toString().length === 1) {
            ms = "0" + m.toString();
        }

        return hs + ":" + ms;
    }
    const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

    const sessionsPerDay = 2;

    const data1 = [
        {date: 1693548000000, time: [8 * 60, 14 * 60]},
        {date: 1693548000000, time: [15 * 60, 17.5 * 60]},
        {date: 1693670400000, time: [6 * 60, 16 * 60]},
        {date: 1693670400000, time: [17 * 60, 20 * 60]},
        {date: 1693717200000, time: [7 * 60, 10 * 60]},
        {date: 1693717200000, time: [14 * 60, 16 * 60]},
    ];

    const data2 = [
        {date: 1693548000000, time: [4 * 60, 10 * 60]},
        {date: 1693548000000, time: [14 * 60, 20 * 60]},
        {date: 1693670400000, time: [9 * 60, 20 * 60]},
        {},
        /*{date: 1693670400000, time: [21 * 60, 22 * 60]},*/
        {date: 1693717200000, time: [5 * 60, 15 * 60]},
        {date: 1693717200000, time: [16 * 60, 19 * 60]},
    ];

    let chElement: HTMLCanvasElement;
    let detailsElement: HTMLDialogElement;

    const init = async () => {
        new Chart(chElement, {
            type: "bar",
            data: {
                labels: data1.map(row => row.date),
                datasets: [
                    {
                        label: "fa",
                        backgroundColor: "blue",
                        data: data1.map(row => row.time),
                    },
                    {
                        label: "ja",
                        backgroundColor: "red",
                        data: data2.map(row => row.time)
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: true
                    },
                    title: {
                        display: false
                    },
                    tooltip: {
                        yAlign: "bottom",
                        callbacks: {
                            title: (items) => {
                                //@ts-ignore
                                let session: number = (items[0]["raw"][1] - items[0]["raw"][0]) / 60;
                                let today = 0;

                                for (let i = 0; i < items[0]["dataset"]["data"].length; i++) {
                                    if (!items[0]["dataset"]["data"][i]) {
                                        continue;
                                    }

                                    //@ts-ignore
                                    if (items[0]["dataset"]["data"][i][0] === items[0]["raw"][0] && items[0]["dataset"]["data"][i][1] === items[0]["raw"][1]) {

                                        if (i % sessionsPerDay) {
                                            //@ts-ignore
                                            if (items[0]["dataset"]["data"][i - 1]) today = (items[0]["dataset"]["data"][i - 1][1] - items[0]["dataset"]["data"][i - 1][0]) / 60
                                        } else {
                                            //@ts-ignore
                                            if (items[0]["dataset"]["data"][i + 1]) today = (items[0]["dataset"]["data"][i + 1][1] - items[0]["dataset"]["data"][i + 1][0]) / 60
                                        }

                                        today = today + session;
                                        break;
                                    }
                                }

                                return "User: " + items[0]["dataset"]["label"] + "\nSession: " + session.toString() + "h\nToday: " + today.toString() + "h";
                            },
                            label: (item) => {
                                //@ts-ignore
                                return formatTime(item["raw"][0]) + "-" + formatTime(item["raw"][1]);
                            },
                        }
                    }
                },
                scales: {
                    x: {
                        type: "timeseries",
                        stacked: false,
                        time: {
                            unit: "day"
                        },
                        grid: {
                            lineWidth: 0.5,
                        }
                    },
                    y: {
                        type: "linear",
                        stacked: false,
                        min: 0,
                        max: 1440,
                        ticks: {
                            stepSize: 60,
                            callback: (value) => {
                                return formatTime((value as number));
                            }
                        },
                        beginAtZero: true,
                    }
                }
            },
        });
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
        <dialog id="details" aria-hidden="true" bind:this={detailsElement} on:click={(e) => {
            if (!e.target) {
                return;
            }

            //@ts-ignore
            if (e.target.id !== "details") {
                return;
            }
             //@ts-ignore
            let width = e.target.clientWidth;
             //@ts-ignore
            let height = e.target.clientHeight;
             //@ts-ignore
            let left = e.target.offsetLeft;
             //@ts-ignore
            let top = e.target.offsetTop;
             //@ts-ignore

            if ((e.x < left || e.x > left + width) || (e.y < top || e.y > top + height)) {
                 //@ts-ignore
                e.target.close();
            }
        }}>
            <table>
                <tr>
                    <th>User</th>
                    <th>Month hours</th>
                    <th>Average hours</th>
                    <th>Month sessions</th>
                    <th>Average sessions</th>
                    <th>File share</th>
                </tr>
                <tr>
                    <td>fa</td>
                    <td>30h</td>
                    <td>8h</td>
                    <td>42</td>
                    <td>2</td>
                    <td>\\PC1\c$</td>
                </tr>
                <tr>
                    <td>ja</td>
                    <td>20h</td>
                    <td>6.5h</td>
                    <td>37</td>
                    <td>1.5</td>
                    <td>\\PC2\c$</td>
                </tr>
            </table>
        </dialog>
        <button on:click={() => detailsElement.showModal()}>details</button>
        <canvas id="ch" bind:this={chElement}></canvas>
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

    #wrapper button {
        position: absolute;
        right: 10px;
        background-color: #3b60dc;
        color: #FFF;
        padding: 3px 5px;
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
