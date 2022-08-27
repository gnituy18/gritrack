import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const { user } = await parent();
  throw redirect(302, `/${user.id}/`);
};
