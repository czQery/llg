<script lang="ts">
    //@ts-ignore
    import Svelecte, {addFormatter} from "svelecte";
    import {getPaletteColor, sleep} from "../ts/helper";
    import type {userInput} from "../ts/global";
    import {userInputStore} from "../ts/global";

    export let userList: userInput[] = [];
    let userListSelect: userInput[] = [];
    let init: boolean = true;

    userInputStore.subscribe((value: userInput[]) => {
        if (value && init) {
            init = false;
            userListSelect = value;
        }
    });

    const change = () => {
        sleep(100);
        userInputStore.set(userListSelect);
    }

    const userColor = (item: userInput, isSelected: boolean) => {
        if (isSelected) {
            for (const u of userListSelect) {
                if (item.text == u.text) {
                    item.color = u.color;
                    break;
                }
            }

            if (!item.color) {
                let color = getPaletteColor();
                let ttl = 0;
                let uIndex = 0;

                for (let i = 0; i < userListSelect.length; i++) {
                    if (userListSelect[i].text == item.text) {
                        uIndex = i;
                    } else if (userListSelect[i].color == color) {
                        color = getPaletteColor();

                        if (ttl > userListSelect.length * userListSelect.length) {
                            continue;
                        }

                        ttl = ttl + 1;
                        i = -1;
                    }
                }

                userListSelect[uIndex].color = color;
                item.color = color;
            }

            return '<div style="border-left: 4px solid ' + item.color + '; padding: 0 0 0 4px">' + item.text + '</div>';
        }

        return '<div>' + item.text + '</div>';
    }
    addFormatter("user-color", userColor);
</script>

<Svelecte inputId="users"
          renderer="user-color"
          options={userList}
          multiple={true}
          collapseSelection={true}
          alwaysCollapsed={true}
          searchable={true}
          valueAsObject={true}
          bind:value={userListSelect}
          on:change={change}/>