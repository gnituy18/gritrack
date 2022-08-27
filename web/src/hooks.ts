import type { Handle, HandleError, GetSession } from "@sveltejs/kit";
import v1 from "$lib/apis/v1";

function isPublicPage(path: string): boolean {
  return path === "/login" || path === "/product";
}

export const handle: Handle = async ({ event, resolve }) => {
  let isLoggedIn = false;

  const cookieStr = event.request.headers.get("cookie");
  if (cookieStr) {
    const sessionIdCookie = cookieStr
      .split(";")
      .find((c) => c.trim().startsWith("sessionid="));
    if (sessionIdCookie) {
      const sessionId = sessionIdCookie.trim().split("=")[1];
      const apiRes = await fetch(v1("/user/current"), {
        headers: { cookie: event.request.headers.get("cookie") },
      });
      if (apiRes.ok) {
        isLoggedIn = true;
        event.locals = { currentUser: await apiRes.json(), sessionId };
      }
    }
  }

  if (!isLoggedIn && !isPublicPage(event.url.pathname)) {
    return new Response(null, {
      status: 302,
      headers: {
        location: "/product",
      },
    });
  } else if (
    isLoggedIn &&
    (event.url.pathname === "/login" || event.url.pathname === "/")
  ) {
    return new Response(null, {
      status: 302,
      headers: {
        location: `/${event.locals.currentUser.id}`,
      },
    });
  }

  return resolve(event);
};

export const handleError: HandleError = async ({ error }) => {
  console.error(error);
};

export const getSession: GetSession = async ({ locals }) => {
  return locals;
};
