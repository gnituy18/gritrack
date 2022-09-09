import { writable } from "svelte/store";
import type { Step } from "$lib/types";
import v1 from "$lib/apis/v1";

class StepsStore {
  private steps = writable<Array<Step>>([]);
  private missionId: string;
  private currentOffset: number = 0;
  private more: boolean;

  public subscribe = this.steps.subscribe;

  public set(missionId: string, steps: Array<Step>, more: boolean, currentOffset: number) {
    this.missionId = missionId;
    this.currentOffset = currentOffset;
    this.steps.set(steps);
    this.more = more;
  }

  public async setRange(missionId: string, offset: number, limit: number) {
    const res = await fetch(
      v1(`/mission/${this.missionId}/step?offset=${offset}&limit=${limit}`),
      {
        credentials: "include",
      }
    );
    const resp = await res.json();
    this.missionId = missionId;
    this.currentOffset = offset + resp.steps.length;
    this.steps.set(resp.steps);
    this.more = resp.more;
  }

  public async updateMore(count: number = 10) {
    const res = await fetch(
      v1(
        `/mission/${this.missionId}/step?offset=${this.currentOffset}&limit=${count}`
      ),
      {
        credentials: "include",
      }
    );
    const moreSteps = await res.json();
    console.log(moreSteps);
    this.currentOffset += moreSteps.steps.length;
    this.more = moreSteps.more;
    this.steps.update((steps) => [...steps, ...moreSteps.steps]);
  }

  public hasMore(): boolean {
    console.log(this.more)
    return this.more;
  }
}

export default new StepsStore();
