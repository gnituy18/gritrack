declare namespace App {
  type User = import("$types").User;

  interface Locals {
    currentUser?: User;
    sessionId?: string;
  }

  interface Platform {}

  interface Session {
    currentUser?: User;
    sessionId?: string;
  }

  interface Stuff {}

  interface PublicEnv {
    PUBLIC_GOOGLE_CLIENT_ID: string;
    PUBLIC_BACKEND_HOST: string;
  }
}
