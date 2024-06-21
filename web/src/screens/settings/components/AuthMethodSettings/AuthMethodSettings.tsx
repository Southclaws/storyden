import { Unready } from "src/components/site/Unready";

import { Heading } from "@/components/ui/heading";

import { SettingsSection } from "../SettingsSection/SettingsSection";

import { Devices } from "./Devices/Devices";
import { OAuth } from "./OAuth/OAuth";
import { Password } from "./Password/Password";
import { useAuthMethodSettings } from "./useAuthMethodSettings";

export function AuthMethodSettings() {
  const state = useAuthMethodSettings();

  if (!state.ready) return <Unready {...state.error} />;

  const { active, available } = state;

  return (
    <SettingsSection gap="4">
      <Heading size="md">Authentication methods</Heading>
      <p>
        We recommend you add more than one method of authentication to your
        account. This will help you recover your account if you lose access to
        one of your devices.
      </p>
      {available.password && <Password active={active.password.length > 0} />}
      {available.webauthn && <Devices active={active.webauthn} />}
      <OAuth active={active.methods} available={available.oauth} />
    </SettingsSection>
  );
}
