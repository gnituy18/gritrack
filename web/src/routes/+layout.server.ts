import type { User } from "$lib/types";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad<{
  currentUser?: User;
  sessionId?: string;
}> = ({ locals }) => {
  const { currentUser, sessionId } = locals;
  return { currentUser, sessionId };
};
