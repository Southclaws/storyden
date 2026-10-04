"use client";

import { parseAsString, parseAsStringLiteral, useQueryState } from "nuqs";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { useAdminSettingsGet } from "@/api/openapi-client/admin";
import { useRoleList } from "@/api/openapi-client/roles";
import { Permission } from "@/api/openapi-schema";
import { AgentRegistrationApprovals } from "@/components/admin/AgentRegistrations/AgentRegistrationApprovals";
import { RegistrationAccessTokens } from "@/components/admin/AgentRegistrations/RegistrationAccessTokens";
import { BackAction } from "@/components/site/Action/Back";
import { Unready } from "@/components/site/Unready";
import { PageHeader } from "@/components/ui/page-header";
import * as Tabs from "@/components/ui/tabs";
import { Text } from "@/components/ui/text";
import { isMemberRole } from "@/lib/role/defaults";
import { useSettings } from "@/lib/settings/settings-client";
import { LStack } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

export function AgentRegistrationsScreen() {
  const [verificationCode] = useQueryState("verification_code", parseAsString);
  const [tab, setTab] = useQueryState(
    "tab",
    parseAsStringLiteral(["tokens", "reviews"]),
  );
  const account = useAccountGet();
  const authorised = hasPermission(account.data, Permission.ADMINISTRATOR);
  const settings = useAdminSettingsGet({ swr: { enabled: authorised } });
  const publicSettings = useSettings();
  const roles = useRoleList({ swr: { enabled: authorised } });

  if (!account.data) return <Unready error={account.error} />;
  if (!authorised)
    return (
      <Unready error="Only administrators can manage agent registrations." />
    );
  if (!settings.data || !publicSettings.ready)
    return <Unready error={settings.error ?? publicSettings.error} />;

  const configuration = settings.data.services?.oauth;
  const oauthEnabled = publicSettings.settings.capabilities.includes("oauth");
  const approvalsEnabled =
    oauthEnabled &&
    configuration?.dynamic_registration_enabled === true &&
    configuration.autonomous_registration_mode === "approval";
  const roleID = configuration?.autonomous_registration_role_id;
  const role = roles.data?.roles.find((role) => role.id === roleID);
  const defaultRoleName =
    !roleID || isMemberRole({ id: roleID })
      ? "Member only"
      : role
        ? `${role.name} + Member`
        : roles.data
          ? "Member only (saved role was deleted)"
          : roles.error
            ? "Role could not be loaded"
            : "Loading role…";

  return (
    <LStack gap="6">
      <PageHeader
        title="Agent registrations"
        back={<BackAction fallbackHref="/admin/authentication" />}
      />
      {!oauthEnabled && (
        <Text variant="supporting">
          Enable OAuth on the server to allow agent registration.
        </Text>
      )}
      {oauthEnabled &&
        (!configuration?.dynamic_registration_enabled ||
          configuration.autonomous_registration_mode === "disabled") && (
          <Text variant="supporting">
            Agent registration is disabled. Tokens can be prepared now and used
            once registration is enabled.
          </Text>
        )}
      <Tabs.Root
        width="full"
        variant="line"
        lazyMount
        value={tab ?? (verificationCode ? "reviews" : "tokens")}
        onValueChange={({ value }) => {
          void setTab(value === "reviews" ? "reviews" : "tokens");
        }}
      >
        <Tabs.List>
          <Tabs.Trigger value="tokens">Registration tokens</Tabs.Trigger>
          <Tabs.Trigger value="reviews">Review registrations</Tabs.Trigger>
          <Tabs.Indicator />
        </Tabs.List>
        <Tabs.Content value="tokens">
          <RegistrationAccessTokens enabled={oauthEnabled} />
        </Tabs.Content>
        <Tabs.Content value="reviews">
          <AgentRegistrationApprovals
            enabled={approvalsEnabled}
            defaultRoleName={defaultRoleName}
          />
        </Tabs.Content>
      </Tabs.Root>
    </LStack>
  );
}
