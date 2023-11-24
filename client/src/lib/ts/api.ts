
import {isJSON, isOK} from "./helper";
import {dataStore, infoStore} from "./global";

export interface infoSum {
    build: string
    selected_users: number
    users: string[] | null
    devices: string[] | null
}

export const loadInfo = async (): Promise<boolean> => {
    const response = await fetch(import.meta.env.VITE_APP_API + "api/info", {
        credentials: "include"
    });

    if (!isJSON(response)) {
        return false;
    }

    const data = await response.json();

    if (!isOK(data)) {
        return false;
    }

    infoStore.set(data["data"] as infoSum);
    return true;
}

export interface dataSum {
    dates: number[]
    items: dataItem[]
}

export interface dataItem {
    name: string
    type: "user" | "device"
    sessions: dataItemSession[]
}

export interface dataItemSession {
    date: number | undefined
    time: number[] | undefined
    detail: string | undefined
}

export const loadData = async (param: string) => {
    const response = await fetch(import.meta.env.VITE_APP_API + "api/data"+param, {
        credentials: "include"
    });

    if (!isJSON(response)) {
        return;
    }

    const data = await response.json();

    dataStore.set(data["data"] as dataSum);
}