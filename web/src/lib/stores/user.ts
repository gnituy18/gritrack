import { writable } from 'svelte/store'

export interface User {
	id: string
	name: string
	email: string
	imageURL: string
}

const user = writable<User>()

export default user
