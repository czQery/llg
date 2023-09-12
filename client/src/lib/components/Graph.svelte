<script lang="ts">
    import "chartjs-adapter-date-fns";
    import {
        BarController,
        BarElement,
        Chart,
        Legend,
        LinearScale,
        PointElement,
        TimeSeriesScale,
        Tooltip
    } from "chart.js";
    import {formatDuration, formatTime, getRandomColor} from "../ts/helper";
    import {type dataSum, type dataUser, type dataUserSession} from "../ts/api";
    import {dataStore, sessionsSums, type sessionSum} from "../ts/global";

    export let aspElement: HTMLSpanElement;
    let chElement: HTMLCanvasElement;
    let data: dataSum = {};

    Chart.register(LinearScale, TimeSeriesScale, BarController, BarElement, PointElement, Tooltip, Legend);
    let chart: Chart<"bar", (number[] | undefined)[], number> | undefined;
    const initChart = (): Chart<"bar", (number[] | undefined)[], number> => {
        return new Chart(chElement, {
            type: "bar",
            data: {
                labels: [],
                datasets: []
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: true,
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
                                let session: number = items[0]["raw"][1] - items[0]["raw"][0];
                                let today = 0;

                                let startIndex = items[0]["dataIndex"] - (items[0]["dataIndex"] % data.sessionsPerDay);

                                for (let i = 0; i < data.sessionsPerDay; i++) {
                                    if (!items[0]["dataset"]["data"][startIndex + i]) continue

                                    //@ts-ignore
                                    today = today + (items[0]["dataset"]["data"][startIndex + i][1] - items[0]["dataset"]["data"][startIndex + i][0]);
                                }

                                let raw = data.users.find((u: dataUser) => u.name == items[0]["dataset"]["label"])?.sessions[items[0]["dataIndex"]]
                                let device: string = "";
                                if (raw && raw.device) {
                                    aspElement.innerText = "\\\\" + raw.device.split(" ")[0] + "\\c$";
                                    device = "\nPC: " + raw.device.split(" ")[0] + "\nIP: " + raw.device.split(" ")[1].replace("(", "").replace(")", "") + "\n";
                                }

                                return "User: " + items[0]["dataset"]["label"] + device + "\nSession: " + formatDuration(session) + "\nToday: " + formatDuration(today);
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
                        //@ts-ignore
                        parsing: false,
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

    dataStore.subscribe(async (value: dataSum) => {
        if (value) {
            data = value;

            if (!chart) {
                chart = initChart();
            }

            chart.data.labels = data.dates.map((d: number) => {
                return d * 1000;
            });

            /*let last: number;
            chart.data.labels = data.users[0].sessions.map((row: dataUserSession) => {
                if (row.date) {
                    last = row.date * 1000
                    return last;
                } else {
                    return last;
                }
            });*/

            let uSums: sessionSum[] = [];

            for (const u of data.users) {
                let color: string = getRandomColor();

                chart.data.datasets.push({
                    label: u.name,
                    backgroundColor: color,
                    data: u.sessions.map((row: dataUserSession) => row.time)
                });

                let uSum = 0;
                for (const s of u.sessions) {
                    if (!s.time) continue;

                    uSum = uSum + (s.time[1] - s.time[0]);
                }

                uSums.push({color: color, sum: uSum});
            }

            chart.update();

            sessionsSums.set(uSums);
        }
    });
</script>

<canvas id="ch" bind:this={chElement}></canvas>

