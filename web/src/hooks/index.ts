import { sequence } from "@sveltejs/kit/hooks";
import type { Handle, HandleError, GetSession, ExternalFetch } from "@sveltejs/kit";
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

  return await resolve(event);
};

export const handle: Handle = sequence(initSession, getUser);

export const handleError: HandleError = async ({ error }) => {
  console.error(error);
};

export const getSession: GetSession = async ({ locals }) => {
  return locals;
};

export const externalFetch: ExternalFetch = async (request) => {
  const cookiesStr = request.headers.get("cookie");

  console.log(request.headers)

  if (!cookiesStr) {
    return fetch(request);
  }
  const cookies = cookiesStr
    .split(";")
    .map((str) => str.trim().split("="))
    .reduce((acc, [key, value]) => ({ ...acc, [key]: value }), {});

  const sessionId = cookies["sessionid"];
  if (sessionId) {
    request.headers.set("sessionid", sessionId);
  }

  console.log(request.headers)

  return fetch(new Request(request.url, request));
};
