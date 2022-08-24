import type { LayoutLoad } from "./__layout";
import v1 from "$/apis/v1";

export const load: LayoutLoad = async ({ fetch, session }) => {
  const res = await fetch(v1("/mission"), {
    headers: {
      sessionid: session.sessionId,
    },
  });

  if (res.status !== 200) {
    throw new Error("@migration task: Migrate this return statement (https://github.com/sveltejs/kit/discussions/5774#discussioncomment-3292693)");
    return {
      status: res.status,
    };
  }

  const missions = await res.json();

  return {
  missions,
};
};
