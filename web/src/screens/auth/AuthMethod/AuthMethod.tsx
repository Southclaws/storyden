import { Spinner } from "@chakra-ui/react";
import { SignUp } from "./Password/SignUp";
import { WebAuthn } from "./WebAuthn/WebAuthn";

interface Props {
  method: string | undefined;
}
export function AuthMethod({ method }: Props) {
  if (!method) {
    return <Spinner />;
  }

  switch (method) {
    case "password":
      return <SignUp />;

    case "webauthn":
      return <WebAuthn />;
  }

  // TODO: Status/error component.
  return <Spinner />;
}
