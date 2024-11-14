import "server-only";

import { AuthProvider } from "src/api/openapi-schema";

import { authProviderList } from "@/api/openapi-server/auth";
import { groupAuthProviders } from "@/lib/auth/utils";

interface Providers {
  password: boolean;
  phone: boolean;
  webauthn: boolean;
  oauth: AuthProvider[];
}

/**
 * Gets available auth providers but separates out Password and Phone providers
 * because they're handled differently. All other providers are pretty standard
 * OAuth2 style providers with links that go off-platform for callbacks etc.
 * @returns Available auth providers with password/phone separated.
 */
export async function getProviders(): Promise<Providers> {
  const { data } = await authProviderList();

  return groupAuthProviders(data.providers);
}
