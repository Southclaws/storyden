import { NextRequest } from "next/server";
import "server-only";

import { AuthProvider } from "src/api/openapi-schema";

import {
  authProviderList,
  getAuthProviderListUrl,
} from "@/api/openapi-server/auth";
import { fetcher } from "@/api/server";
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
  const { data } = await authProviderList(
    {
      revalidate: 0,
      cache: "no-store",
    } as any /* HACK: "revalidate" is passed to next options. */,
  );

  return groupAuthProviders(data.providers);
}
