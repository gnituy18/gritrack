import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const { currentUser } = await parent();
  throw redirect(302, `/${currentUser.id}/`);
};
