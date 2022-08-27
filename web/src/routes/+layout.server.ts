import type { User } from "$lib/types";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad<{
  user?: User;
  sessionId?: string;
}> = ({ locals }) => {
  const { user, sessionId } = locals;
  return { user, sessionId };
};
