<script lang="ts">
    //@ts-ignore
    import Svelecte, {addFormatter} from "svelecte";
    import {getPaletteColor, sleep} from "../ts/helper";
    import type {itemInput} from "../ts/global";
    import {itemInputStore} from "../ts/global";

    export let op: object = [];
    let itemListSelect: itemInput[] = [];
    let init: boolean = true;

    itemInputStore.subscribe((value: itemInput[]) => {
        if (value && init) {
            init = false;
            itemListSelect = value;
        }
    });

    const change = () => {
        sleep(100);
        itemInputStore.set(itemListSelect);
    }

    const itemColor = (item: itemInput, isSelected: boolean) => {
        if (isSelected) {
            for (const u of itemListSelect) {
                if (item.value == u.value) {
                    item.color = u.color;
                    break;
                }
            }

            if (!item.color) {
                let color = getPaletteColor();
                let ttl = 0;
                let uIndex = 0;

                for (let i = 0; i < itemListSelect.length; i++) {
                    if (itemListSelect[i].value == item.value) {
                        uIndex = i;
                    } else if (itemListSelect[i].color == color) {
                        color = getPaletteColor();

                        if (ttl > itemListSelect.length * itemListSelect.length) {
                            continue;
                        }

                        ttl = ttl + 1;
                        i = -1;
                    }
                }

                itemListSelect[uIndex].color = color;
                item.color = color;
            }

            return '<div style="border-left: 4px solid ' + item.color + '; padding: 0 0 0 4px">' + item.value + '</div>';
        }

        return '<div>' + item.value + '</div>';
    }
    addFormatter("item-color", itemColor);
</script>

<Svelecte inputId="item"
          renderer="item-color"
          groupLabelField="groupHeader"
          groupItemsField="items"
          options={op}
          multiple={true}
          collapseSelection={true}
          alwaysCollapsed={true}
          searchable={true}
          valueAsObject={false}
          bind:readSelection={itemListSelect}
          on:change={change}/>