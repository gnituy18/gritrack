declare namespace App {
  type User = import("$lib/types").User;

  interface Locals {
    currentUser?: User;
    sessionId?: string;
  }

  interface Platform {}

  interface Stuff {}

  interface PublicEnv {
    PUBLIC_GOOGLE_CLIENT_ID: string;
    PUBLIC_BACKEND_HOST: string;
  }
}
