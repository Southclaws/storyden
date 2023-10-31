import { Heading } from "@chakra-ui/react";

import { Unready } from "src/components/site/Unready";

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
      {available.password && <Password />}
      {available.webauthn && <Devices active={active.webauthn} />}
      <OAuth active={active.methods} available={available.oauth} />
    </SettingsSection>
  );
}
