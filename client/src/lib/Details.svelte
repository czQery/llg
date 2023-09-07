<script lang="ts">
    let detailsElement: HTMLDialogElement;
</script>

<div>
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
    <button id="details-btn" on:click={() => detailsElement.showModal()}>details</button>
</div>

<style>
    #details-btn {
        position: absolute;
        right: 10px;
        background-color: #3b60dc;
        color: #FFF;
        padding: 3px 5px;
    }
</style>