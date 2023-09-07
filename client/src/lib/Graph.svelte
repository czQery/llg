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
    import {onMount} from "svelte";
    import {formatTime} from "../helper";

    Chart.register(LinearScale, TimeSeriesScale, BarController, BarElement, PointElement, Tooltip, Legend);

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

    onMount(() => {
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
    });
</script>

<canvas id="ch" bind:this={chElement}></canvas>

