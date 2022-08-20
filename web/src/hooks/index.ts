import type { Handle, HandleError, GetSession } from "@sveltejs/kit";
import v1 from "$apis/v1";

function isPublicPage(path: string): boolean {
  return path === "/" || path === "/login";
}

export const handle: Handle = async ({ event, resolve }) => {
  let isLoggedIn = false;

  const cookieStr = event.request.headers.get("cookie");
  if (cookieStr) {
    const sessionIdCookie = cookieStr.split(";").find((c) => c.trim().startsWith("sessionid="));
    if (sessionIdCookie) {
      const sessionId = sessionIdCookie.trim().split("=")[1];
      const apiRes = await fetch(v1("/user/current"), {
        headers: { cookie: event.request.headers.get("cookie") },
      });
      if (apiRes.ok) {
        isLoggedIn = true;
        event.locals = { currentUser: await apiRes.json(), sessionId }
      }
    }
  }

  if (!isLoggedIn || !isPublicPage(event.url.pathname)) {
  }


  if (!sessionId && !isPublicPage(event.url.pathname)) {
    return new Response(null, {
      status: 302,
      headers: {
        location: "/product",
      },
    });
  }


  if (apiRes.status === 401) {
    return new Response(null, {
      status: 302,
      headers: {
        location: "/product",
      },
    });
  }

  if (!apiRes.ok && !isPublicPage(event.url.pathname)) {
    return new Response(null, {
      status: 302,
      headers: {
        location: "/product",
      },
    });
  }

  if (apiRes.ok && event.url.pathname === "/login") {
    return new Response(null, {
      status: 302,
      headers: {
        location: `/${event.locals.currentUser.id}`,
      },
    });
  }

  event.locals = {
    ...event.locals,
    currentUser: await apiRes.json(),
  };

  return resolve(event);
});

export const handleError: HandleError = async ({ error }) => {
  console.error(error);
};

export const getSession: GetSession = async ({ locals }) => {
  return locals;
};
