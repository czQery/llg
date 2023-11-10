import {writable, type Writable} from "svelte/store";
import type {dataSum, infoSum} from "./api";

// infoStore: data from /api/info
export const infoStore: Writable<infoSum> = writable();

// dataStore: data from /api/data
export const dataStore: Writable<dataSum> = writable();

export interface userActive {
    name: string
    color: string
    sum: number
}

// userActiveStore: list of rendered users
export const userActiveStore: Writable<userActive[]> = writable();

export interface userInput {
    value: string
    color: string
}

// userInputStore: list selected users
export const userInputStore: Writable<userInput[]> = writable();