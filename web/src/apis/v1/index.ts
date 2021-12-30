export default function (url: string): string {
  return `${import.meta.env.ENV_BACKEND_HOST}/api/v1${url}`;
}
