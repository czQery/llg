import {type Writable, writable} from "svelte/store";
import {isJSON, isOK} from "./helper";

export interface dataSum {
    dates: number[]
    users: dataUser[]
    sessionsPerDay: number
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

export const dataStore: Writable<dataSum> = writable();
export const loadData = async () => {
    const response = await fetch(import.meta.env.VITE_APP_API + "api/data", {
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