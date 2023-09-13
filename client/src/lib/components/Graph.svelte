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
            if (!value.dates || !value.users) {
                return;
            }

            data = value;

            if (!chart) {
                chart = initChart();
            } else {
                chart.data.datasets = [];
                if (chart.tooltip) {
                    chart.tooltip.setActiveElements([], {x: 0, y: 0});
                }
            }

            chart.data.labels = data.dates.map((d: number) => {
                return d * 1000;
            });

            let uSums: sessionSum[] = [];

            for (let i = 0; i < data.users.length; i++) {
                if (!data.users[i].sessions) {
                    continue;
                }

                let color: string = getRandomColor();
                let hidden = false;

                if (i > 2) {
                    hidden = true;
                }

                chart.data.datasets.push({
                    label: data.users[i].name,
                    backgroundColor: color,
                    data: data.users[i].sessions.map((row: dataUserSession) => row.time),
                    hidden: hidden
                });

                let uSum = 0;
                for (const s of data.users[i].sessions) {
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

