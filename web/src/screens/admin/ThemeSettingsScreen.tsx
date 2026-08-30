"use client";

import { useState } from "react";
import { toast } from "sonner";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { adminThemeDelete, useAdminThemeGet } from "@/api/openapi-client/admin";
import { Permission } from "@/api/openapi-schema";
import { UnreadyBanner } from "@/components/site/Unready";
import { revalidateTheme } from "@/lib/theme/revalidate-theme";
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
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string>();

  if (!account.data && !account.error) return <UnreadyBanner />;
  if (!hasPermission(account.data, Permission.ADMINISTRATOR)) {
    return (
      <UnreadyBanner error="Administrator permission is required to manage custom themes." />
    );
  }
  if (!query.data) {
    return <UnreadyBanner error={query.error ?? account.error} />;
  }

  async function disableTheme() {
    setBusy(true);
    setError(undefined);
    try {
      await adminThemeDelete();
      await Promise.all([query.mutate(), revalidateTheme()]);
      setThemeEditingEnabled(false);
      toast.success("Custom theme disabled");
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause));
    } finally {
      setBusy(false);
    }
  }

  return (
    <ThemeSettingsView
      assets={[...query.data.stylesheets, ...query.data.scripts]}
      active={query.data.enabled}
      runtimeDisabled={query.data.runtime_disabled}
      editingEnabled={editingEnabled}
      busy={busy}
      error={error}
      onEnableEditing={() => setThemeEditingEnabled(true)}
      onExitEditing={() => setThemeEditingEnabled(false)}
      onDisable={disableTheme}
    />
  );
}
