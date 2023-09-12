import {writable, type Writable} from "svelte/store";
import type {dataSum} from "./api";

// dataStore: data from /api/data
export const dataStore: Writable<dataSum> = writable();



export interface sessionSum {
    color: string
    sum: number
}

// sessionsSums: users month hour sum
export const sessionsSums: Writable<sessionSum[]> = writable();