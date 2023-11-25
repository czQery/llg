import {type itemInput} from "./global";
import {getPaletteColor} from "./helper";
import {type infoSum} from "./api";

export const parseSelection = (type: "user" | "device", urlParams: URLSearchParams, info: infoSum): itemInput[] => {
    let url: string | null = urlParams.get(type+"s");

    let data: itemInput[] = [];
    if (url) {
        let urlList: string[] = url.split(",").filter((u: string): boolean => {

            if (type == "user" && info.users) {
                for (const uInfo of info.users) {
                    if (uInfo == u) {
                        return true;
                    }
                }
            }

            if (type == "device" && info.devices) {
                for (const uInfo of info.devices) {
                    if (uInfo == u) {
                        return true;
                    }
                }
            }

            return false;
        })

        data = urlList.map((u: string): itemInput => {
            return <itemInput>{type: type, value: u, color: ""};
        });
    } else {
        if (info.users && type === "user") {
            data = info.users.slice(0, info.selected_items).map((u: string): itemInput => {
                return <itemInput>{type: type, value: u, color: ""}
            });
        }
    }

    for (const u of data) {
        u.color = getPaletteColor();
    }

    return data;
}