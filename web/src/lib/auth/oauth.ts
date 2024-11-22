import { filter } from "lodash/fp";

import { AuthProvider } from "@/api/openapi-schema";

const hasLink = (provider: AuthProvider): provider is OAuthProvider => {
  return provider.link !== undefined;
};

export const filterWithLink = (list: AuthProvider[]): OAuthProvider[] =>
  filter(hasLink)(list);

export type OAuthProvider = AuthProvider & { link: string };
