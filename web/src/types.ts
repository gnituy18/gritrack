export type Mission = {
  id: string;
  name: string;
  description: string;
};

export type Step = {
  id: string;
  summary: string;
  items: Array<Item>;
  createdAt: number;
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
