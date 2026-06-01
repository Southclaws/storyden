import { WithDisclosure } from "src/utils/useDisclosure";

import { useAccountView } from "@/api/openapi-client/accounts";
import { InstanceCapability, ProfileReference } from "@/api/openapi-schema";
import { Unready } from "@/components/site/Unready";
import { useSettings } from "@/lib/settings/settings-client";

import { MemberPasswordResetForm } from "./MemberPasswordResetForm";

export type Props = {
  profile: ProfileReference;
};

export function MemberPasswordResetDialog(props: WithDisclosure<Props>) {
  const accountResult = useAccountView(props.profile.id);
  const settingsResult = useSettings();

  if (!settingsResult.ready) {
    return <Unready error={settingsResult.error} />;
  }

  if (!accountResult.data) {
    return <Unready error={accountResult.error} />;
  }

  const hasEmail = settingsResult.settings.capabilities.includes(
    InstanceCapability.email_client,
  );

  return (
    <MemberPasswordResetForm
      profile={props.profile}
      account={accountResult.data}
      hasEmail={hasEmail}
      onClose={props.onClose}
    />
  );
}
