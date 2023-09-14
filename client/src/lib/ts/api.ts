
import {isJSON, isOK} from "./helper";
import {dataStore} from "./global";

export interface dataSum {
    dates: number[]
    users: dataUser[]
}

export interface dataUser {
    name: string
    sessions: dataUserSession[]
}

export interface dataUserSession {
    date: number | undefined
    device: string | undefined
    time: number[] | undefined
}
export const loadData = async (param: string) => {
    const response = await fetch(import.meta.env.VITE_APP_API + "api/data"+param, {
        credentials: "include"
    });

    if (!isJSON(response)) {
        return;
    }

    const data = await response.json();

    if (!isOK(data)) {
        return;
    }

    dataStore.set(data["data"] as dataSum);
}