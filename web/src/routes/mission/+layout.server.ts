import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ parent }) => {
  const { sessionId, user } = await parent();
  return { sessionId, user };
};
