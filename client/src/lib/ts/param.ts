import {type userInput} from "./global";
import {getPaletteColor} from "./helper";
import {type infoSum} from "./api";

export const parseSelection = (name: string, urlParams: URLSearchParams, info: infoSum): userInput[] => {
    let url: string = urlParams.get(name);
    let data: userInput[] = [];
    if (url) {
        let urlList: string[] = url.split(",").filter((u: string): boolean => {

            if (name == "users" && info.users) {
                for (const uInfo: string of info.users) {
                    if (uInfo == u) {
                        return true;
                    }
                }
            }

            if (name == "devices" && info.devices) {
                for (const uInfo: string of info.devices) {
                    if (uInfo == u) {
                        return true;
                    }
                }
            }

            return false;
        })

        data = urlList.map((u: string): userInput => {
            return {value: u, color: ""};
        });
    } else {
        if (info.users) {
            data = info.users.slice(0, info.selected_users).map((u: string): userInput => {
                return {value: u, color: ""}
            });
        }
    }

    for (const u: userInput of data) {
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
}