import {
  useOAuthClientList,
  useOAuthRefreshTokenList,
} from "@/api/openapi-client/auth";
import { OAuthTokenSettings } from "@/components/settings/OAuthTokenSettings/OAuthTokenSettings";
import { Unready } from "@/components/site/Unready";

export function MemberOAuthSettingsScreen() {
  const tokens = useOAuthRefreshTokenList();
  const clients = useOAuthClientList();

  if (!tokens.data || !clients.data) {
    return <Unready error={tokens.error ?? clients.error} />;
  }

  return (
    <OAuthTokenSettings
      tokens={tokens.data.tokens}
      clients={clients.data.clients}
    />
  );
}
