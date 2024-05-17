import { Unready } from "src/components/site/Unready";

import { SettingsSection } from "../SettingsSection/SettingsSection";

import { Heading2 } from "@/components/ui/typography-heading";

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
      <Heading2 size="md">Authentication methods</Heading2>
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
