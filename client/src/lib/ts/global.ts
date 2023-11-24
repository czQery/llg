import {writable, type Writable} from "svelte/store";
import type {dataSum, infoSum} from "./api";

// infoStore: data from /api/info
export const infoStore: Writable<infoSum> = writable();

// dataStore: data from /api/data
export const dataStore: Writable<dataSum> = writable();

export interface itemActive {
    name: string
    color: string
    sum: number
}

// itemActiveStore: list of rendered items
export const itemActiveStore: Writable<itemActive[]> = writable();

export interface itemInput {
    type: "user" | "device"
    value: string
    color: string
}

// itemInputStore: list selected items
export const itemInputStore: Writable<itemInput[]> = writable();