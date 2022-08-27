import type { LayoutLoad } from "./$types";
import v1 from "$lib/apis/v1";

export const load: LayoutLoad = async ({ fetch, parent }) => {
  const { sessionId } = await parent();
  const res = await fetch(v1("/mission"), {
    headers: {
      sessionid: sessionId,
    },
  });

  if (res.status !== 200) {
    return {
      status: res.status,
    };
  }

  const missions = await res.json();

  return {
    missions,
  };
};
