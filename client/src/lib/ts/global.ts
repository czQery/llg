import {writable, type Writable} from "svelte/store";
import type {dataSum, infoSum} from "./api";

// infoStore: data from /api/info
export const infoStore: Writable<infoSum> = writable();

// dataStore: data from /api/data
export const dataStore: Writable<dataSum> = writable();

export interface sessionSum {
    color: string
    sum: number
}

// sessionSumStore: users month hour sum
export const sessionSumStore: Writable<sessionSum[]> = writable();