declare namespace App {
  type User = import("$lib/types").User;

  interface Locals {
    user?: User;
  }

  interface PageData {
    user?: User;
  }

  interface PublicEnv {
    PUBLIC_GOOGLE_CLIENT_ID: string;
    PUBLIC_BACKEND_HOST: string;
  }
}
