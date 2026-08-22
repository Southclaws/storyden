import { filter } from "lodash/fp";

import { AuthProvider } from "@/api/openapi-schema";

const hasLink = (provider: AuthProvider): provider is OAuthProvider => {
  return provider.link !== undefined;
};

export const filterWithLink = (list: AuthProvider[]): OAuthProvider[] =>
  filter(hasLink)(list);

export type OAuthProvider = AuthProvider & { link: string };

export function formatOAuthGrant(grant: string): string {
  if (grant === "client_credentials") return "client credentials";
  if (grant === "authorization_code") return "authorization code";
  if (grant === "refresh_token") return "refresh token";
  if (grant.includes("device_code")) return "device code";

  return grant;
}
