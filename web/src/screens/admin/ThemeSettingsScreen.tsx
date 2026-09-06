"use client";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { useAdminThemeGet } from "@/api/openapi-client/admin";
import { Permission } from "@/api/openapi-schema";
import { UnreadyBanner } from "@/components/site/Unready";
import {
  setThemeEditingEnabled,
  useThemeEditingEnabled,
} from "@/lib/theme/theme-editing";
import { hasPermission } from "@/utils/permissions";

import { ThemeSettingsView } from "./ThemeSettingsView";

export function ThemeSettingsScreen() {
  const account = useAccountGet();
  const query = useAdminThemeGet({
    swr: { enabled: hasPermission(account.data, Permission.ADMINISTRATOR) },
  });
  const editingEnabled = useThemeEditingEnabled();

  if (!account.data && !account.error) return <UnreadyBanner />;
  if (!hasPermission(account.data, Permission.ADMINISTRATOR)) {
    return (
      <UnreadyBanner error="Administrator permission is required to manage custom themes." />
    );
  }
  if (!query.data) {
    return <UnreadyBanner error={query.error ?? account.error} />;
  }

  return (
    <ThemeSettingsView
      assets={[...query.data.stylesheets, ...query.data.scripts]}
      editingEnabled={editingEnabled}
      onEnableEditing={() => setThemeEditingEnabled(true)}
      onExitEditing={() => setThemeEditingEnabled(false)}
    />
  );
}
