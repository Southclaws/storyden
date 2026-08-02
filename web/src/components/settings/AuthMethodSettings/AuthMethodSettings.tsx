"use client";

import { Unready } from "@/components/site/Unready";
import { PageHeading } from "@/components/ui/page-heading";
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
      <LStack>
        <PageHeading>Authentication methods</PageHeading>
        <p>
          We recommend you add more than one method of authentication to your
          account. This will help you recover your account if you lose access to
          one of your devices.
        </p>
      </LStack>

      {available.password && <Password active={active.password.length > 0} />}

      {/* NOTE: WebAuthn is not enabled as a 2FA yet. */}
      {/* {available.webauthn && <Devices active={active.webauthn} />} */}

      {(available.oauth.length > 0 || active.methods.length > 0) && (
        <OAuth active={active.methods} available={available.oauth} />
      )}
    </LStack>
  );
}
