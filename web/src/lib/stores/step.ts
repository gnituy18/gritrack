import { writable } from "svelte/store";
import type { Step } from "$lib/types";
import v1 from "$lib/apis/v1";

class StepsStore {
  private steps = writable<Array<Step>>([]);
  private missionId: string;
  private currentOffset: number = 0;
  private more: boolean ;

  public subscribe = this.steps.subscribe;

  public set(missionId: string, steps: Array<Step>, currentOffset: number) {
    this.missionId = missionId;
    this.currentOffset = currentOffset;
    this.steps.set(steps);
  }

  public async setRange(missionId: string, offset: number, limit: number) {
    const res = await fetch(
      v1(`/mission/${this.missionId}/step?offset=${offset}&limit=${limit}`),
      {
        credentials: "include",
      }
    );
    const steps = await res.json();
    this.missionId = missionId;
    this.currentOffset = offset + steps.length;
    this.steps.set(steps);
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
    this.currentOffset += moreSteps.length;
    if (moreSteps.length < count) {
      this.more = false;
    }
    this.steps.update((steps) => [...steps, ...moreSteps]);
  }

  public hasMore(): boolean {
    return this.more;
  }
}

export default new StepsStore();
