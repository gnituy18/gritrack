export interface Mission {
	id: string;
	name: string;
	description: string;
}

export interface Step {
	id: string;
	summary: string;
	items: Array<Item>;
}

export interface Item {
	type: number;
	desc: string;

	time?: ItemTime;
}

export interface ItemTime {
	duration: number;
}
