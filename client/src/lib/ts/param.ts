import {type itemInput} from "./global";
import {getPaletteColor} from "./helper";
import {type infoSum} from "./api";

export const parseSelection = (name: "user" | "device", urlParams: URLSearchParams, info: infoSum): itemInput[] => {
    let url: string | null = urlParams.get(name+"s");

    let data: itemInput[] = [];
    if (url) {
        let urlList: string[] = url.split(",").filter((u: string): boolean => {

            if (name == "user" && info.users) {
                for (const uInfo of info.users) {
                    if (uInfo == u) {
                        return true;
                    }
                }
            }

            if (name == "device" && info.devices) {
                for (const uInfo of info.devices) {
                    if (uInfo == u) {
                        return true;
                    }
                }
            }

            return false;
        })

        data = urlList.map((u: string): itemInput => {
            return <itemInput>{type: name, value: u, color: ""};
        });
    } else {
        if (info.users) {
            data = info.users.slice(0, info.selected_users).map((u: string): itemInput => {
                return <itemInput>{type: name, value: u, color: ""}
            });
        }
    }

    for (const u of data) {
        let color: string = getPaletteColor();
        let ttl: number = 0;

        for (let i: number = 0; i < data.length; i++) {
            if (data[i].color == color && data[i].value != u.value) {
                color = getPaletteColor();

                if (ttl > data.length * data.length) {
                    continue;
                }

                ttl = ttl + 1;
                i = -1;
            }
        }
        u.color = color;
    }

    return data;
}