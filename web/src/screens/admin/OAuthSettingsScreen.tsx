"use client";

import { parseAsInteger, useQueryStates } from "nuqs";

import {
  useAdminOAuthClientList,
  useAdminOAuthDeviceAuthorisationList,
  useAdminOAuthRefreshTokenList,
  useAdminSettingsGet,
} from "@/api/openapi-client/admin";
import { OAuthSettings } from "@/components/admin/OAuthSettings/OAuthSettings";
import { Unready } from "@/components/site/Unready";
import { parseAdminSettings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";

export function OAuthSettingsScreen() {
  const [filters] = useQueryStates({
    page: parseAsInteger.withDefault(1),
  });

  const clients = useAdminOAuthClientList();
  const settings = useAdminSettingsGet();
  const publicSettings = useSettings();
  const devices = useAdminOAuthDeviceAuthorisationList();
  const tokens = useAdminOAuthRefreshTokenList({
    page: filters.page.toString(),
  });

  if (!settings.data || !publicSettings.ready) {
    return <Unready error={settings.error ?? publicSettings.error} />;
  }
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
      settings={{
        ...parseAdminSettings(settings.data),
        capabilities: publicSettings.settings.capabilities,
      }}
      clients={clients.data.clients}
      deviceAuthorisations={devices.data.device_authorisations}
      tokens={tokens.data.tokens}
      tokenPage={{
        currentPage: filters.page,
        totalPages: tokens.data.total_pages,
        pageSize: tokens.data.page_size,
      }}
    />
  );
}
