import { Account } from "@/api/openapi-schema";
import { Member } from "@/lib/settings/member-settings";
import { Settings } from "@/lib/settings/settings";
import { useSessionData } from "@/lib/settings/session-client";

export function useSession(
  initial?: Account,
  initialSettings?: Settings,
): Member | undefined {
  const { session } = useSessionData(initial, initialSettings);

  return session;
}

export function useAccountSession(
  initial?: Account,
  initialSettings?: Settings,
) {
  const query = useSessionData(initial, initialSettings);

  return {
    ...query,
    data: query.data?.account,
  };
}
