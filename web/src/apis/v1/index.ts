import { env } from "$env/dynamic/public";

export default function (url: string): string {
  return `${env.PUBLIC_BACKEND_HOST}/api/v1${url}`;
}
