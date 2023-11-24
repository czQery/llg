<script lang="ts">
    import {BarController, BarElement, Chart, LinearScale, PointElement, Tooltip} from "chart.js";
    import {formatDate, formatDuration, formatTime, sleep} from "../ts/helper";
    import {type dataItem, type dataItemSession, type dataSum} from "../ts/api";
    import {dataStore, type itemActive, itemActiveStore, type itemInput} from "../ts/global";

    export let itemInputList: itemInput[] = [];
    export let aspElement: HTMLSpanElement;
    let chElement: HTMLCanvasElement;
    let data: dataSum = {dates: [], items: []};

    Chart.register(LinearScale, BarController, BarElement, PointElement, Tooltip);
    let chart: Chart<"bar", (number[] | undefined)[], number> | undefined;
    const initChart = async (): Promise<Chart<"bar", (number[] | undefined)[], number>> => {

        while (true) {
            if (chElement) {
                break;
            }
            await sleep(50);
        }

        return new Chart(chElement, {
            type: "bar",
            data: {
                labels: [],
                datasets: []
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                maxBarThickness: 80,
                plugins: {
                    legend: {
                        display: false,
                        labels: {
                            font: {
                                size: 15
                            }
                        },
                    },
                    title: {
                        display: false
                    },
                    tooltip: {
                        yAlign: "bottom",
                        //@ts-ignore
                        events: ["click", "touchstart"],
                        titleFont: {
                            size: 15
                        },
                        bodyFont: {
                            size: 15
                        },
                        callbacks: {
                            title: (items) => {
                                //@ts-ignore
                                let session: number = items[0]["raw"]["y"][1] - items[0]["raw"]["y"][0];
                                let day = 0;

                                for (const d of items[0]["dataset"]["data"]) {
                                    //@ts-ignore
                                    if (d["x"] == items[0]["raw"]["x"]) {
                                        //@ts-ignore
                                        day = day + (d["y"][1] - d["y"][0]);
                                    }
                                }

                                let raw = data.items.find((u: dataItem) => u.name == items[0]["dataset"]["label"]);
                                let rawSession = raw?.sessions[items[0]["dataIndex"]];
                                let detail: string = "";

                                if (!rawSession || !rawSession.detail || !raw!.type) {
                                    return "Name: " + items[0]["dataset"]["label"] + detail + "\nSession: " + formatDuration(session) + "\nDay: " + formatDuration(day);
                                }

                                if (raw?.type === "user") {
                                    aspElement.innerText = "\\\\" + rawSession.detail.split(" ")[0] + "\\c$";
                                    detail = "\nPC: " + rawSession.detail.split(" ")[0] + "\nIP: " + rawSession.detail.split(" ")[1].replace("(", "").replace(")", "") + "\nType: user\n";
                                    return "User: " + items[0]["dataset"]["label"] + detail + "\nSession: " + formatDuration(session) + "\nDay: " + formatDuration(day);
                                }

                                if (raw?.type === "device") {
                                    aspElement.innerText = "\\\\" + items[0]["dataset"]["label"]!.split(" ")[0] + "\\c$";
                                    detail = "\nPC: " + items[0]["dataset"]["label"]!.split(" ")[0] + "\nIP: " + items[0]["dataset"]["label"]!.split(" ")[1].replace("(", "").replace(")", "") + "\nType: device\n";
                                    return "User: " + rawSession.detail + detail + "\nSession: " + formatDuration(session) + "\nDay: " + formatDuration(day);
                                }
                            },
                            label: (item) => {
                                //@ts-ignore
                                return [formatTime(item["raw"]["y"][0]) + " - " + formatTime(item["raw"]["y"][1]) + " (" + formatDate(item["raw"]["x"]) + ")"];
                            },
                        }
                    }
                },
                scales: {
                    x: {
                        type: "linear",
                        stacked: false,
                        ticks: {
                            stepSize: 1,
                            callback: (value) => {
                                if ((value as number) < 10) {
                                    return null;
                                }

                                if (data.dates.length == 1) {
                                    let timeThis = new Date((value as number) * 60 * 60 * 24 * 1000);
                                    let timeData = new Date(data.dates[0] * 60 * 60 * 24 * 1000);

                                    if (timeThis.getMonth() != timeData.getMonth() || timeThis.getFullYear() != timeData.getFullYear()) {
                                        return "";
                                    }
                                }

                                return formatDate((value as number));
                            }
                        },
                        beginAtZero: false,
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

    dataStore.subscribe(async (value: dataSum) => {
        if (!chart) {
            chart = await initChart();
        } else {
            chart.data.labels = [];
            chart.data.datasets = [];
            if (chart.tooltip) {
                chart.tooltip.setActiveElements([], {x: 0, y: 0});
            }

            //@ts-ignore
            chart.options.scales.x.min = null;
            //@ts-ignore
            chart.options.scales.x.max = null;

            chart.update();
            itemActiveStore.set([]);
        }

        if (!value || !value.dates || !value.items) {
            return;
        }

        data = value;

        chart.data.labels = data.dates;

        let startDate = new Date(data.dates[0] * 60 * 60 * 24 * 1000);
        let endDate = new Date(data.dates[0] * 60 * 60 * 24 * 1000);

        let rangeMin = new Date(startDate.getFullYear(), startDate.getMonth(), 2);
        let rangeMax = new Date(endDate.getFullYear(), endDate.getMonth() + 1, 1);

        //@ts-ignore
        chart.options.scales.x.min = Math.trunc(rangeMin.valueOf() / 60 / 60 / 24 / 1000);
        //@ts-ignore
        chart.options.scales.x.max = Math.trunc(rangeMax.valueOf() / 60 / 60 / 24 / 1000);

        let aUsers: itemActive[] = [];

        for (let i = 0; i < data.items.length; i++) {

            if (!data.items[i].sessions) {
                continue;
            }

            let color: string = "";
            for (const u of itemInputList) {
                if (u.value.toLowerCase() == data.items[i].name.toLowerCase().split(" ")[0]) {
                    color = u.color;
                    break;
                }
            }

            if (color === "") {
                for (const u of itemInputList) {
                    if (u.value.toLowerCase().includes(data.items[i].name.toLowerCase().split(" ")[0])) {
                        color = u.color;
                        break;
                    }
                }
            }

            chart.data.datasets.push({
                label: data.items[i].name,
                backgroundColor: color,
                //@ts-ignore
                data: data.items[i].sessions.map((row: dataItemSession) => {
                    return {x: row.date, y: row.time}
                }),
                hidden: false
            });

            let uSum = 0;
            for (const s of data.items[i].sessions) {
                if (!s.time) continue;

                uSum = uSum + (s.time[1] - s.time[0]);
            }

            aUsers.push({name: data.items[i].name, color: color, sum: uSum});
        }

        chart.update();
        itemActiveStore.set(aUsers);
    });
</script>

<canvas id="ch" bind:this={chElement}></canvas>

