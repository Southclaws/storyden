import { keyBy, values } from "lodash/fp";
import "server-only";

import {
  AuthProvider,
  AuthProviderListOKResponse,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";

interface Providers {
  password: boolean;
  phone: boolean;
  webauthn: boolean;
  providers: AuthProvider[];
}

const group = keyBy<AuthProvider>("provider");

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

  // pull out password and phone, if present, the rest are OAuth2 providers.
  const { password, phone, webauthn, ...rest } = group(providers);

  return {
    password: Boolean(password),
    phone: Boolean(phone),
    webauthn: Boolean(webauthn),
    providers: values(rest),
  };
}
