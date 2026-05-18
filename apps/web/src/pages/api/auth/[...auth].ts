import { Auth } from "@auth/core";
import GitHub from "@auth/core/providers/github";
import Credentials from "@auth/core/providers/credentials";

export const prerender = false;

const authConfig = {
  providers: [
    GitHub({
      clientId: import.meta.env.GITHUB_CLIENT_ID,
      clientSecret: import.meta.env.GITHUB_CLIENT_SECRET,
    }),
    Credentials({
      name: "Email",
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        // TODO: Validate credentials against SurrealDB user records.
        // This placeholder keeps the Auth.js route wired while the user schema is integrated.
        if (!credentials?.email) {
          return null;
        }

        return {
          id: String(credentials.email),
          email: String(credentials.email),
          name: String(credentials.email).split("@")[0],
        };
      },
    }),
  ],
  secret: import.meta.env.AUTH_SECRET,
  trustHost: true,
};

export const ALL = ({ request }: { request: Request }) => Auth(request, authConfig);
