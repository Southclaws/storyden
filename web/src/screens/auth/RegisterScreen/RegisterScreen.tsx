import { AuthMode, RegistrationMode } from "@/api/openapi-schema";
import { authProviderList } from "@/api/openapi-server/auth";
import { tServer } from "@/i18n/server";
import { VStack } from "@/styled-system/jsx";

import { RegisterEmailForm } from "./RegisterEmail/RegisterEmailForm";
import { RegisterHandleForm } from "./RegisterHandle/RegisterHandleForm";
import { RegisterPhoneForm } from "./RegisterPhone/RegisterPhoneForm";

type Props = {
  invitationID?: string;
  registrationMode: RegistrationMode;
};

export async function RegisterScreen({
  invitationID,
  registrationMode,
}: Props) {
  const { data } = await authProviderList({
    cache: "no-store",
  });

  const isInviteOnly = registrationMode === RegistrationMode.invitation;
  if (isInviteOnly) {
    return (
      <VStack textAlign="center">
        <styled.h1 fontWeight="bold">Registration is invite-only.</styled.h1>
        <styled.p color="fg.muted" textWrap="balance">
          Ask a community member or administrator for an invitation link to
          join.
        </styled.p>
      </VStack>
    );
  }

  const isDisabled = registrationMode === RegistrationMode.disabled;
  if (isDisabled) {
    return (
      <VStack textAlign="center">
        <styled.h1 fontWeight="bold">
          Registration is currently closed.
        </styled.h1>
        <styled.p color="fg.muted" textWrap="balance">
          This site has closed public registration of accounts.
        </styled.p>
      </VStack>
    );
  }

  switch (data.mode) {
    case AuthMode.handle:
      return (
        <RegisterHandleForm webauthn={false} invitationID={invitationID} />
      );

    case AuthMode.email:
      return <RegisterEmailForm invitationID={invitationID} />;

    case AuthMode.phone:
      return <RegisterPhoneForm />;

    default:
      console.error("no authentication modes available");

      return (
        <VStack>
          <p>{await tServer("This instance is private.")}</p>
        </VStack>
      );
  }
}
