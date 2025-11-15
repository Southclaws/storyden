import { useAccountGet } from "@/api/openapi-client/accounts";
import { Account } from "@/api/openapi-schema";
import { Member, parseMemberSettings } from "@/lib/settings/member-settings";
import { Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";

export function useSession(
  initial?: Account,
  initialSettings?: Settings,
): Member | undefined {
  const { settings } = useSettings(initialSettings);
  const { data } = useAccountGet({
    swr: {
      fallbackData: initial,
    },
  });

  if (!data || !settings) {
    return undefined;
  }

  const session = parseMemberSettings(settings.metadata, data);

  return session;
}
