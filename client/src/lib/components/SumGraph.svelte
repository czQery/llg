<script lang="ts">
    import {type userActive, userActiveStore} from "../ts/global";
    import {formatDuration} from "../ts/helper";

    let sums: userActive[] = [];
    let sumsSum: number = 0;

    userActiveStore.subscribe((value: userActive[]) => {
        if (value) {
            sums = value;
            sumsSum = 0;

            for (const s of sums) {
                sumsSum = sumsSum + s.sum;
            }
        }
    });
</script>

<div id="sm">
    {#if sumsSum <= 0}
        <div style="color: #000; width: 100%">There are no data for selected date or users!</div>
    {:else}
        {#each sums as s}
            {#if s.sum !== -1}
                <div style={"background-color: "+s.color+"; width: "+(s.sum / sumsSum) * 100+"%"}>{((window.innerWidth - (2 * 30)) * (s.sum / sumsSum) < 40) ? "" : formatDuration(s.sum)}</div>
            {/if}
        {:else}
            <div style="color: #000; width: 100%">There are no data for selected date or users!</div>
        {/each}
    {/if}
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