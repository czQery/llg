<script lang="ts">
    import {BarController, BarElement, Chart, LinearScale, PointElement, Tooltip} from "chart.js";
    import {formatDate, formatDuration, formatTime, sleep} from "../ts/helper";
    import {type dataSum, type dataUser, type dataUserSession} from "../ts/api";
    import {type userActive, userActiveStore, dataStore, type userInput} from "../ts/global";

    export let userInputList: userInput[] = [];
    export let aspElement: HTMLSpanElement;
    let chElement: HTMLCanvasElement;
    let data: dataSum = {dates: [], users: []};

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

                                let raw = data.users.find((u: dataUser) => u.name == items[0]["dataset"]["label"])?.sessions[items[0]["dataIndex"]]
                                let device: string = "";
                                if (raw && raw.device) {
                                    aspElement.innerText = "\\\\" + raw.device.split(" ")[0] + "\\c$";
                                    device = "\nPC: " + raw.device.split(" ")[0] + "\nIP: " + raw.device.split(" ")[1].replace("(", "").replace(")", "") + "\n";
                                }

                                return "User: " + items[0]["dataset"]["label"] + device + "\nSession: " + formatDuration(session) + "\nDay: " + formatDuration(day);
                            },
                            label: (item) => {
                                //@ts-ignore
                                return [formatTime(item["raw"]["y"][0]) + " - " + formatTime(item["raw"]["y"][1]) + " (" + formatDate(item["raw"]["x"])+")"];
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
            userActiveStore.set([]);
        }

        if (!value || !value.dates || !value.users) {
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

        let aUsers: userActive[] = [];

        for (let i = 0; i < data.users.length; i++) {
            if (!data.users[i].sessions) {
                continue;
            }

            let color: string = "";
            for (const u of userInputList) {
                if (u.text.toLowerCase() == data.users[i].name.toLowerCase()) {
                    color = u.color;
                    break;
                }
            }

            chart.data.datasets.push({
                label: data.users[i].name,
                backgroundColor: color,
                //@ts-ignore
                data: data.users[i].sessions.map((row: dataUserSession) => {
                    return {x: row.date, y: row.time}
                }),
                hidden: false
            });

            let uSum = 0;
            for (const s of data.users[i].sessions) {
                if (!s.time) continue;

                uSum = uSum + (s.time[1] - s.time[0]);
            }

            aUsers.push({name: data.users[i].name, color: color, sum: uSum});
        }

        chart.update();
        userActiveStore.set(aUsers);
    });
</script>

<canvas id="ch" bind:this={chElement}></canvas>

