import type { PageLoad } from "./$types";
import v1 from "$lib/apis/v1";
import type { Steps } from "$lib/types";

export const load: PageLoad = async ({ params, fetch, parent }) => {
  const { user } = await parent();
  try {
    const userId = params.userId;
    const missionName = params.missionName;
    let res = await fetch(v1(`/user/${userId}/missionName/${missionName}`), {
      headers: {
        sessionid: user.sessionId,
      },
    });
    if (!res.ok) {
      return {
        status: res.status,
      };
    }

    const mission = await res.json();
    res = await fetch(v1(`/mission/${mission.id}/step?offset=0&limit=10`), {
      credentials: "include",
    });

    if (!res.ok) {
      return {
        status: res.status,
      };
    }
    const resp: Steps = await res.json();

    return {
      mission,
      steps: resp.steps,
      more: resp.more,
      currentOffset: resp.steps.length,
    };
  } catch (error) {
    console.error(error);
    throw error(500);
  }
};
