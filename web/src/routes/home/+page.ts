import { redirect } from '@sveltejs/kit';
import type { PageLoad } from "../__layout";

export const load: PageLoad = async ({ session }) => {
  throw redirect(302, `/${session.currentUser.id}/`);
};
