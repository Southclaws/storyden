import { useAccessKeyList } from "@/api/openapi-client/auth";
import { AccessKeysSettings } from "@/components/settings/AccessKeysSettings/AccessKeysSettings";
import { Unready } from "@/components/site/Unready";

export function MemberAccessKeysSettingsScreen() {
  const { data, error } = useAccessKeyList();
  if (!data) {
    return <Unready error={error} />;
  }

  return <AccessKeysSettings keys={data.keys} />;
}
