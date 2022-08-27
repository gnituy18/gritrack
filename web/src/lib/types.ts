export type User = {
  id: string;
  sessionId: string;
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
  // TODO refactor this
  date: string;
  time?: number;
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
