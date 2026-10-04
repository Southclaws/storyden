"use client";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { useAdminSettingsGet } from "@/api/openapi-client/admin";
import { useRoleList } from "@/api/openapi-client/roles";
import { Permission } from "@/api/openapi-schema";
import { AuthenticationSettingsForm } from "@/components/admin/AuthenticationSettings/AuthenticationSettingsForm";
import { UnreadyBanner } from "@/components/site/Unready";
import { parseAdminSettings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";
import { hasPermission } from "@/utils/permissions";

export function AuthenticationSettingsScreen() {
  const settings = useAdminSettingsGet();
  const publicSettings = useSettings();
  const account = useAccountGet();
  const roles = useRoleList();
  if (!settings.data || !publicSettings.ready || !roles.data || !account.data) {
    return (
      <UnreadyBanner
        error={
          settings.error ?? publicSettings.error ?? roles.error ?? account.error
        }
      />
    );
  }

  return (
    <AuthenticationSettingsForm
      settings={{
        ...parseAdminSettings(settings.data),
        capabilities: publicSettings.settings.capabilities,
      }}
      roles={roles.data.roles}
      canManageRoles={hasPermission(account.data, Permission.MANAGE_ROLES)}
      canApproveRegistrations={hasPermission(
        account.data,
        Permission.ADMINISTRATOR,
      )}
    />
  );
}
