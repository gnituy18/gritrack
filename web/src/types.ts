export type User = {
  id: string;
  name: string;
  email: string;
  picture: string;
};

export type Mission = {
  id: string;
  name: string;
  description: string;
};

export type Step = {
  id: string;
  createdAt: number;
  summary: string;
  items: Array<Item>;
};

export enum ItemType {
  Time = 1,
}

export type Item = {
  type: ItemType.Time;
  desc: string;
  time: ItemTime;
};

export type ItemTime = {
  duration: number;
};
