export type User = {
  id: string;
  name: string;
  email: string;
  picture: string;
};

export type Mission = {
  id: string;
  userId: string;
  name: string;
  description: string;
};

export type Step = {
  id: string;
  createdAt: number;
  summary: string;
  items: Array<Item>;
  date: string;
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

export type DropdownItem = {
  label: string;
  action: (label: string) => void;
};
