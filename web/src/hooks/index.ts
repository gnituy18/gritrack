import { sequence } from '@sveltejs/kit/hooks';
import type { Handle } from '@sveltejs/kit';

const initSession: Handle = async ({ request, resolve }) => {
	request.locals = {
		session: {},
	}
	return resolve(request)
}

const getUser: Handle = async ({ request, resolve }) => {
	const apiRes = await fetch('http://localhost:8080/api/v1/user/current', {
		headers: { ...request.headers },
	});

	if (!apiRes.ok && request.path !== '/login') {
		return {
			status: 302,
			headers: {
				location: '/login',
			},
		}
	}

	if (apiRes.ok && request.path === '/login') {
		return {
			status: 302,
			headers: {
				location: '/',
			},
		}
	}

	if (apiRes.ok) {
		request.locals = {
			...request.locals,
			session: {
				...request.locals.session,
				currentUser: await apiRes.json(),
			},
		}
	}

	const res = await resolve(request)
	return {
		...res,
		headers: {
			...res.headers,
			'set-cookie': apiRes.headers.get('set-cookie'),
		},
	}
}

export const handle: Handle<{
	currentUser: {
		picture: string;
		name: string;
	}
}> = sequence(initSession, getUser)

export async function getSession({ locals: { session } }) {
	return session
}
