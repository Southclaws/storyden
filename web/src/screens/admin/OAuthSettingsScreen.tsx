import {
  useAdminOAuthClientList,
  useAdminOAuthDeviceAuthorisationList,
  useAdminOAuthRefreshTokenList,
} from "@/api/openapi-client/admin";
import { OAuthSettings } from "@/components/admin/OAuthSettings/OAuthSettings";
import { Unready } from "@/components/site/Unready";

export function OAuthSettingsScreen() {
  const clients = useAdminOAuthClientList();
  const devices = useAdminOAuthDeviceAuthorisationList();
  const tokens = useAdminOAuthRefreshTokenList();

  if (!clients.data) {
    return <Unready error={clients.error} />;
  }
  if (!devices.data) {
    return <Unready error={devices.error} />;
  }
  if (!tokens.data) {
    return <Unready error={tokens.error} />;
  }

  return (
    <OAuthSettings
      clients={clients.data.clients}
      deviceAuthorisations={devices.data.device_authorisations}
      tokens={tokens.data.tokens}
    />
  );
}
