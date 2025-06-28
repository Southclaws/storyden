import { useAdminAccessKeyList } from "@/api/openapi-client/admin";
import { AccessKeySettings } from "@/components/admin/AccessKeySettings/AccessKeySettings";
import { UnreadyBanner } from "@/components/site/Unready";

export function AccessKeySettingsScreen() {
  const { data, error } = useAdminAccessKeyList();
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return <AccessKeySettings keys={data.keys} />;
}
