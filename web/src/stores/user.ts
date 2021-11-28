import { writable } from "svelte/store";

export interface User {
  id: string;
  name: string;
  email: string;
  picture: string;
}

const user = writable<User>();

export default user;
