import { sequence } from "@sveltejs/kit/hooks";
import type { Handle, HandleError, GetSession } from "@sveltejs/kit";
import v1 from "$apis/v1";

const initSession: Handle = async ({ event, resolve }) => {
  event.locals = {
    currentUser: null,
  };
  return resolve(event);
};

const getUser: Handle = async ({ event, resolve }) => {
  const apiRes = await fetch(v1("/user/current"), {
    headers: { cookie: event.request.headers.get("cookie") },
  });

  if (!apiRes.ok && event.url.pathname !== "/login") {
    return new Response(null, { status: 302, headers: { location: "/login" } });
  }

  if (apiRes.ok && event.url.pathname === "/login") {
    return new Response(null, { status: 302, headers: { location: "/" } });
  }

  if (apiRes.ok) {
    event.locals = {
      ...event.locals,
      currentUser: await apiRes.json(),
    };
  }

  const res = await resolve(event);
  res.headers.set("set-cookie", apiRes.headers.get("set-cookie"));
  return res;
};

export const handle: Handle = sequence(initSession, getUser);

export const handleError: HandleError = async ({ error }) => {
  console.error(error);
};

export const getSession: GetSession = async ({ locals }) => {
  return locals;
};
