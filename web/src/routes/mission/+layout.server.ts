import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ parent }) => {
  const { sessionId, currentUser } = await parent();
  return { sessionId, currentUser };
};
