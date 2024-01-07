<script lang="ts">
    //@ts-ignore
    import Svelecte, {addFormatter} from "svelecte";
    import {getPaletteColor, sleep} from "../ts/helper";
    import type {itemInput} from "../ts/global";
    import {itemInputStore} from "../ts/global";
    import {IconUser, IconPcCase, IconHelp} from "../ts/icon";

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
        let style: string = "--bg: var(--c0)";

        if (isSelected) {
            for (const u of itemListSelect) {
                if (item.value == u.value) {
                    item.color = u.color;
                    break;
                }
            }

            if (!item.color) {
                let color = getPaletteColor();
                let uIndex = 0;

                for (let i = 0; i < itemListSelect.length; i++) {
                    if (itemListSelect[i].value == item.value) {
                        uIndex = i;
                    }
                }

                itemListSelect[uIndex].color = color;
                item.color = color;
            }

            style = "--bg: "+item.color;

            //return '<div style="border-left: 4px solid ' + item.color + '; padding: 0 0 0 4px">' + item.value + '</div>';
        }

        switch (item.type) {
            case "user":
                return '<div style="'+style+'">' + IconUser + item.value + '</div>';
            case "device":
                return '<div style="'+style+'">' + IconPcCase + item.value + '</div>';
            default:
                return '<div style="'+style+'">' + IconHelp + item.value + '</div>';
        }

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
          valueAsObject={true}
          highlightFirstItem={true}
          placeholder="Select"
          bind:value={itemListSelect}
          on:change={change}/>