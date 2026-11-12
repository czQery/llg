<script lang="ts">
    import Svelecte, { addRenderer } from "svelecte"
    import { getPaletteColor, sleep } from "../ts/helper"
    import type { itemInput } from "../ts/global"
    import { itemInputStore } from "../ts/global"
    import { IconX, IconHelp, IconPcCase, IconUser } from "../ts/icon"

    interface Props {
        op?: object
    }

    let { op = [] }: Props = $props()
    let itemListSelect: itemInput[] = $state([])
    let init: boolean = true;

    itemInputStore.subscribe((value: itemInput[]) => {
        if (value && init) {
            init = false
            itemListSelect = value;
        }
    })

    const change = () => {
        sleep(100)
        itemInputStore.set(itemListSelect)
    }

    const itemColor = (itemRaw: object, isSelected?: boolean, inputValue?: string) => {
        let item = itemRaw as itemInput;
        let style: string = "--bg: var(--c0)";

        if (isSelected) {
        	console.log("selected", item.value, item.color)
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
                return '<div style="' + style + '">' + IconUser + item.value + "</div>";
            case "device":
                return '<div style="' + style + '">' + IconPcCase + item.value + "</div>";
            default:
                return '<div style="' + style + '">' + IconHelp + item.value + "</div>";
        }
    };

    addRenderer("item-color", itemColor);
</script>

{#snippet listHeader()}
	{#if itemListSelect.length !== 0}
	    <div class="sv-dropdown--selection">
	        {#each itemListSelect as item (item.value)}
	            <div class="sv-item--wrap in-dropdown">
	                <div class="sv-item--content">
	                    {@html itemColor(item, true)}
	                </div>
	                <button class="sv-item--btn" data-action="deselect" data-id={item.value} tabindex="-1" type="button" aria-label="Remove {item.value}">
	                	{@html IconX}
	                </button>
	            </div>
	        {/each}
	    </div>
    {/if}
{/snippet}

<Svelecte
    inputId="items"
    renderer="item-color"
    groupLabelField="groupHeader"
    groupItemsField="items"
    options={op}
    multiple={true}
    collapseSelection="always"
    keepSelectionInList={false}
    resetOnBlur={true}
    searchable={true}
    valueField="value"
    labelField="value"
    valueAsObject={true}
    highlightFirstItem={true}
    placeholder="Select"
    bind:value={itemListSelect}
    onChange={change}
    {listHeader}
/>

<style>
	:global(.sv-control) {
        border: 1px solid #FFF !important;
        border-radius: var(--rad) !important;
        min-height: 30px !important;
        height: 30px !important;
        color: #000 !important;
        outline: none !important;
        height: 30px;
        width: 230px !important;
        flex-grow: 0 !important;
    }

    :global(.sv-control:hover) {
        cursor: text;
    }

    :global(.svelecte.is-focused .sv-control) {
        border: 1px solid var(--cb) !important;
        outline: none !important;
    }

    :global(.indicator-icon) {
        transition: var(--transition) color linear !important;
    }

    :global(.indicator-icon:hover) {
        cursor: pointer !important;
    }

    :global(.sv-input--text) {
        font-size: 16px !important;
        color: #000 !important;
    }

    :global(.sv-input--text:focus::placeholder) {
        color: #000 !important;
        opacity: 0 !important;
    }

    :global(.sv-input--text::placeholder) {
        color: #000 !important;
        opacity: 0.4 !important;
    }

    :global(.sv-item--wrap) {
        margin: 2px !important;
        border-radius: var(--rad) !important;
    }

    :global(.sv-item--btn) {
        border-radius: var(--rad) !important;
        transition: var(--transition) all linear !important;
    }

    :global(.sv-item--btn:hover) {
        background-color: var(--sv-dropdown-active-bg);
        filter: var(--hover);
    }

    :global(.sv-item--btn:active) {
        animation: reverse click-scale 100ms !important;
    }

    :global(.sv-item--btn svg) {
    	height: 16px;
      	width: 16px;
        fill: var(--c0) !important;
    }

    :global(.sv-btn-indicator) {
    	margin: 9px !important;
    }

    :global(.sv-control--selection) {
    	padding-left: 6px !important;
    }

    :global(.sv_dropdown) {
        border-radius: var(--rad) !important;
        min-width: 230px !important;
        width: 230px !important;
    }

    :global(.sv-dropdown--selection) {
    	padding: 4px;
    }

    :global(.sv-dropdown--selection > .sv-item--wrap) {
	    background-color: var(--sv-dropdown-active-bg);
    }

    :global(.sv-dropdown--selection > .sv-item--wrap > .sv-item--content) {
    	width: 100%
    }

    :global(.sv-dropdown--selection > .sv-item--wrap > .sv-item--btn) {
		position: absolute;
		display: grid;
		align-items: center;
		box-sizing: border-box;
		height: 25px;
		right: 0;
		width: 25px;
		justify-items: center;
		margin: -3px 7px 0 0;
    }

    :global(.sv-dropdown-scroll .sv-optgroup-header) {
        color: #000;
        text-align: left !important;
        cursor: default !important;
    }

    :global(.sv-item--content) {
        border-radius: var(--rad) !important;
        color: #000 !important;
        cursor: pointer !important;
    }

    :global(.sv-item--content div) {
        display: flex;
        width: 100%;
        max-width: 183px;
        overflow: hidden;
    }

    :global(.sv-item--wrap.in-dropdown) {
        transition: var(--transition) all linear !important;
    }

    :global(.sv-item--wrap.in-dropdown .sv-item--content div) {
        max-width: 192px !important;
    }

    :global(.sv-item--content div svg) {
        height: 19px;
        width: 19px;
        margin: 0 6px 0 0;
        padding: 2px 0;
        border-radius: var(--rad);
        background-color: var(--bg);
        stroke: white;
        display: inline-block;
        min-width: 19px;
    }
</style>
