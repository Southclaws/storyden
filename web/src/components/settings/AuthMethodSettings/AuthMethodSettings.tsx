"use client";

import { Unready } from "@/components/site/Unready";
import { PageHeading } from "@/components/ui/page-heading";
import { Text } from "@/components/ui/text";
import { LStack } from "@/styled-system/jsx";

import { OAuth } from "./OAuth/OAuth";
import { Password } from "./Password/Password";
import { useAuthMethodSettings } from "./useAuthMethodSettings";

export function AuthMethodSettings() {
  const { ready, error, data } = useAuthMethodSettings();
  if (!ready) {
    return <Unready error={error} />;
  }

  const { active, available } = data;

  return (
    <LStack gap="4">
      <LStack gap="1">
        <PageHeading>Authentication methods</PageHeading>
        <Text variant="supporting">Manage how you log in to this site.</Text>
      </LStack>

      {available.password && <Password active={active.password.length > 0} />}
      {(available.oauth.length > 0 || active.methods.length > 0) && (
        <OAuth active={active.methods} available={available.oauth} />
      )}
    </LStack>
  );
}
