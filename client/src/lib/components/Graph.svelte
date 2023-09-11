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
    import {formatTime, getRandomColor} from "../ts/helper";
    import {dataStore, type dataUser, type dataUserSession} from "../ts/api";

    export let aspElement: HTMLSpanElement;
    let chElement: HTMLCanvasElement;
    let data: dataUser[] = [];
    const sessionsPerDay = 4;


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

                                let raw = data.find((u) => u.name == items[0]["dataset"]["label"])?.sessions[items[0]["dataIndex"]]
                                let device: string = "";

                                if (raw && raw.device) {
                                    aspElement.innerText = "\\\\" + raw.device.split(" ")[0] + "\\c$";
                                    device = "\nPC: " + raw.device.split(" ")[0] + "\nIP: " + raw.device.split(" ")[1].replace("(", "").replace(")", "") + "\n";
                                }

                                return "User: " + items[0]["dataset"]["label"] + device + "\nSession: " + session.toFixed(1).toString() + "h\nToday: " + today.toFixed(1).toString() + "h";
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

    dataStore.subscribe(async (value: dataUser[]) => {
        if (value) {
            data = value;

            if (!chart) {
                chart = initChart();
            }

            let last: number;
            chart.data.labels = data[0].sessions.map((row: dataUserSession) => {
                if (row.date) {
                    last = row.date * 1000
                    return last;
                } else {
                    return last;
                }
            });

            for (const u of data) {
                chart.data.datasets.push({
                    label: u.name,
                    backgroundColor: getRandomColor(),
                    data: u.sessions.map((row: dataUserSession) => row.time)
                });
            }

            chart.update();
        }
    });
</script>

<canvas id="ch" bind:this={chElement}></canvas>

