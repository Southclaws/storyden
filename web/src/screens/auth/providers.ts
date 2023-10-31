import "server-only";

import {
  AuthProvider,
  AuthProviderListOKResponse,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { groupAuthProviders } from "src/components/settings/utils";

interface Providers {
  password: boolean;
  phone: boolean;
  webauthn: boolean;
  providers: AuthProvider[];
}

/**
 * Gets available auth providers but separates out Password and Phone providers
 * because they're handled differently. All other providers are pretty standard
 * OAuth2 style providers with links that go off-platform for callbacks etc.
 * @returns Available auth providers with password/phone separated.
 */
export async function getProviders(): Promise<Providers> {
  const { providers } = await server<AuthProviderListOKResponse>({
    url: "/v1/auth",
  });

  return groupAuthProviders(providers);
}
