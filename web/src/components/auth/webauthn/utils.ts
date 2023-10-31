import {
  startAuthentication,
  startRegistration,
} from "@simplewebauthn/browser";
import "client-only";

import {
  webAuthnGetAssertion,
  webAuthnMakeAssertion,
  webAuthnMakeCredential,
  webAuthnRequestCredential,
} from "src/api/openapi/auth";

export async function passkeyLogin(handle: string) {
  const { publicKey } = await webAuthnGetAssertion(handle);
  const credential = await startAuthentication(publicKey);

  // HACK:
  // 1. https://github.com/MasterKale/SimpleWebAuthn/issues/330
  // 2. https://github.com/go-webauthn/webauthn/issues/93
  credential.response.userHandle = undefined;

  await webAuthnMakeAssertion(credential);
}

export async function passkeyRegister(handle: string) {
  const { publicKey } = await webAuthnRequestCredential(handle);

  const credential = await startRegistration({
    ...publicKey,
  });

  await webAuthnMakeCredential(credential);
}
