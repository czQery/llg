<script lang="ts">
    import {sessionSumStore, type sessionSum} from "../ts/global";
    import {formatDuration} from "../ts/helper";

    let sums: sessionSum[] = [];
    let sumsSum: number = 0;

    sessionSumStore.subscribe(async (value: sessionSum[]) => {
        if (value) {
            sums = value;
            sumsSum = 0;
        }

        for (const s of sums) {
            sumsSum = sumsSum + s.sum;
        }
    });
</script>

<div id="sm">
    {#each sums as s}
        <div style={"background-color: "+s.color+"; width: "+(s.sum / sumsSum) * 100+"%"}>{formatDuration((s.sum / sumsSum) * 100 * 60)}</div>
    {/each}
</div>

<style>
    #sm {
        display: flex;
        grid-column: 1/3;
        grid-row: 3/4;
        height: 30px;
        padding: 0 10px 10px 10px;
        background-color: #FFF;
        border-radius: 0 0 var(--rad) var(--rad);
        box-shadow: var(--shadow);
    }

    #sm div {
        height: 20px;
        display: inline-block;
        text-align: center;
        color: white;
        /*border-radius: var(--rad);*/
    }
</style>