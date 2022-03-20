import { writable } from "svelte/store";
import type { Mission, Step } from "$types";

export const missions = writable<Array<Mission>>([]);

export const steps = writable<Array<Step>>([]);
