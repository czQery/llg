<script lang="ts">
    //@ts-ignore
    import Svelecte, {addFormatter} from "svelecte";
    import {getPaletteColor, sleep} from "../ts/helper";
    import type {itemInput} from "../ts/global";
    import {itemInputStore} from "../ts/global";
    import {IconHelp, IconPcCase, IconUser} from "../ts/icon";
    import {onMount} from "svelte";

    export let op: object = [];
    let itemListSelect: itemInput[] = [];
    let init: boolean = true;
    let itemInput: any;

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

            style = "--bg: " + item.color;
        }

        switch (item.type) {
            case "user":
                return '<div style="' + style + '">' + IconUser + item.value + '</div>';
            case "device":
                return '<div style="' + style + '">' + IconPcCase + item.value + '</div>';
            default:
                return '<div style="' + style + '">' + IconHelp + item.value + '</div>';
        }

    }
    addFormatter("item-color", itemColor);

    onMount(() => {
        let input = document.getElementById("items") as HTMLInputElement;
        let control = document.querySelector(".sv-control") as HTMLElement
        if (input && control) {
            input.onmousedown = () => {
                input.blur();
            }
            control.onmousedown = () => {
                if (control.classList.contains("is-active")) itemInput.focus();
            }
            input.onfocus = () => {
                if (control.classList.contains("is-active")) input.blur();
            }
        }
    });
</script>

<Svelecte inputId="items"
          renderer="item-color"
          groupLabelField="groupHeader"
          groupItemsField="items"
          options={op}
          multiple={true}
          collapseSelection={true}
          resetOnBlur={true}
          alwaysCollapsed={true}
          searchable={true}
          valueAsObject={true}
          highlightFirstItem={true}
          placeholder="Select"
          bind:value={itemListSelect}
          bind:this={itemInput}
          on:change={change}/>