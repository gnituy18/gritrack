import { sequence } from "@sveltejs/kit/hooks";
import type { Handle, HandleError, GetSession } from "@sveltejs/kit";
import v1 from "$apis/v1";

const extractSessionId: Handle = async ({ event, resolve }) => {
  const cookieStr = event.request.headers.get("cookie");
  if (!cookieStr) {
    return resolve(event);
  }

  const sessionIdCookie = cookieStr.split(";").find((c) => c.trim().startsWith("sessionid="));

  let sessionId: string;
  if (sessionIdCookie) {
    sessionId = sessionIdCookie.trim().split("=")[1];
  }

  event.locals = {
    sessionId,
  };
  return resolve(event);
};

const getUser: Handle = async ({ event, resolve }) => {
  const { sessionId } = event.locals;
  if (sessionId) {
    const apiRes = await fetch(v1("/user/current"), {
      headers: { cookie: event.request.headers.get("cookie") },
    });

    if (apiRes.ok) {
      event.locals = {
        ...event.locals,
        currentUser: await apiRes.json(),
      };
    }
  }

  return await resolve(event);
};

export const handle: Handle = sequence(extractSessionId, getUser);

export const handleError: HandleError = async ({ error }) => {
  console.error(error);
};

export const getSession: GetSession = async ({ locals }) => {
  return locals;
};
