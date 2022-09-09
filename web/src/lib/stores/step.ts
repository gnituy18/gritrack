import { writable } from "svelte/store";
import type { Step, Steps } from "$lib/types";
import v1 from "$lib/apis/v1";

class StepsStore {
  private steps = writable<Steps>({ steps: [], more: false });
  private missionId: string;
  private currentOffset: number = 0;

  public subscribe = this.steps.subscribe;

  public set(
    missionId: string,
    steps: Array<Step>,
    more: boolean,
    currentOffset: number
  ) {
    this.missionId = missionId;
    this.currentOffset = currentOffset;
    this.steps.set({ steps, more });
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
    this.currentOffset += moreSteps.steps.length;
    this.steps.update((steps) => ({
      more: moreSteps.more,
      steps: [...steps.steps, ...moreSteps.steps],
    }));
  }
}

export default new StepsStore();
