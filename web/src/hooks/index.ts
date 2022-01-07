import { sequence } from "@sveltejs/kit/hooks";
import type { Handle, HandleError } from "@sveltejs/kit";
import v1 from "$apis/v1";

const initSession: Handle = async ({ request, resolve }) => {
  request.locals = {
    session: {},
  };
  return resolve(request);
};

const getUser: Handle = async ({ request, resolve }) => {
  const apiRes = await fetch(v1("/user/current"), {
    headers: { cookie: request.headers.cookie },
  });

  if (!apiRes.ok && request.url.pathname !== "/login") {
    return {
      status: 302,
      headers: {
        location: "/login",
      },
    };
  }

  if (apiRes.ok && request.url.pathname === "/login") {
    return {
      status: 302,
      headers: {
        location: "/",
      },
    };
  }

  if (apiRes.ok) {
    request.locals = {
      ...request.locals,
      session: {
        ...request.locals.session,
        currentUser: await apiRes.json(),
      },
    };
  }

  const res = await resolve(request);
  return {
    ...res,
    headers: {
      ...res.headers,
      "set-cookie": apiRes.headers.get("set-cookie"),
    },
  };
};

export const handle: Handle<{
  currentUser: {
    picture: string;
    name: string;
  };
}> = sequence(initSession, getUser);

export const handleError: HandleError = async ({ error }) => {
  console.error(error);
};

export async function getSession({ locals: { session } }) {
  return session;
}
