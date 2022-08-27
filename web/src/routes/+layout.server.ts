import type { User } from "$lib/types";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad<{ user?: User }> = ({ locals }) => {
  const { user } = locals;
  return { user };
};
