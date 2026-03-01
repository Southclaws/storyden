"use client";

import { useGetSession } from "@/api/openapi-client/misc";
import { Account } from "@/api/openapi-schema";

import { Member, parseMemberSettings } from "./member-settings";
import { Settings, parseSettings } from "./settings";

export type SessionData = {
  settings?: Settings;
  session?: Member;
};

export function useSessionData(
  initialAccount?: Account,
  initialSettings?: Settings,
  revalidateOnMount = false,
) {
  const query = useGetSession({
    swr: {
      fallbackData: initialSettings
        ? {
            info: initialSettings,
            account: initialAccount,
          }
        : undefined,
      revalidateOnMount,
    },
  });

  const settings = query.data ? parseSettings(query.data.info) : undefined;
  const session = parseMemberSettings(query.data?.account, settings?.metadata);

  return {
    ...query,
    settings,
    session,
  } satisfies typeof query & SessionData;
}
