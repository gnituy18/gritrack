import { writable } from "svelte/store";
import type { Mission } from "$lib/types";
import v1 from "$lib/apis/v1";

class MissionsStore {
  private missions = writable<Array<Mission>>([]);

  public set = this.missions.set;
  public subscribe = this.missions.subscribe;

  public async update() {
    const res = await fetch(v1("/mission"), { credentials: "include" });
    const missions = await res.json();
    this.missions.set(missions);
  }

  public async delete(missionId: string) {
    const res = await fetch(v1(`/mission/${missionId}`), {
      method: "DELETE",
      credentials: "include",
    });

    if (!res.ok) {
      throw new Error(await res.text());
    }

    this.missions.update((missions) =>
      missions.filter((mission) => mission.id !== missionId)
    );
  }

  public async create(name: string) {
    const res = await fetch(v1("/mission"), {
      method: "POST",
      credentials: "include",
      headers: {
        "content-type": "application/json",
      },
      body: JSON.stringify({
        name: name,
      }),
    });

    if (!res.ok) {
      throw new Error(await res.text());
    }

    await this.update();
  }
}

export default new MissionsStore();
