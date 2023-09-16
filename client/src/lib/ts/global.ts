import {writable, type Writable} from "svelte/store";
import type {dataSum, infoSum} from "./api";

// infoStore: data from /api/info
export const infoStore: Writable<infoSum> = writable();

// dataStore: data from /api/data
export const dataStore: Writable<dataSum> = writable();

export interface activeUser {
    name: string
    color: string
    sum: number
}

// activeUserStore: list of rendered users
export const activeUserStore: Writable<activeUser[]> = writable();