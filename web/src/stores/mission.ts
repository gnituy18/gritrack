import { writable } from "svelte/store";
import type { Step } from "$types";

export const steps = writable<Array<Step>>([]);
