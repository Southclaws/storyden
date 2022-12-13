import { FieldValues, SubmitHandler, useForm } from "react-hook-form";
import {
  webAuthnMakeCredential,
  webAuthnRequestCredential,
} from "src/api/openapi/auth";

import {
  RegistrationCredential,
  RegistrationCredentialJSON,
} from "@simplewebauthn/typescript-types";

export function useWebAuthn() {
  const { register, handleSubmit, formState } = useForm({
    //
  });

  const onSubmit: SubmitHandler<FieldValues> = async ({
    username,
  }: FieldValues) => {
    const response = await webAuthnRequestCredential(username);

    const { publicKey } = response;
    if (!publicKey) {
      throw new Error("publicKey is empty");
    }

    const publicKeyCredentialCreationOptions = {
      ...publicKey,

      rp: {
        id: "localhost", // TODO: Configure self-domain or get from API.
        name: "Storyden",
      },

      // overwrite challenge and user with the correct format
      challenge: Uint8Array.from(publicKey.challenge as string, (c) =>
        c.charCodeAt(0)
      ),
      user: {
        id: Uint8Array.from(publicKey.user.id as string, (c) =>
          c.charCodeAt(0)
        ),
        name: publicKey.user.name,
        displayName: publicKey.user.displayName,
      },
    };

    const credential = (await navigator.credentials.create({
      publicKey: publicKeyCredentialCreationOptions,
    })) as RegistrationCredential; // cast required as TS is outdated currently.

    if (credential == null) {
      throw new Error("credential was null after create");
    }

    const makeCredentialBody: RegistrationCredentialJSON = {
      id: credential?.id,
      rawId: bufToBase64(credential.rawId),
      type: credential.type,
      response: {
        clientDataJSON: bufToBase64(credential.response.clientDataJSON),
        attestationObject: bufToBase64(credential.response.attestationObject),
      },
      clientExtensionResults: credential.getClientExtensionResults(),
      authenticatorAttachment: credential.authenticatorAttachment,
    };

    console.log({ credential, makeCredentialBody });

    const creds = await webAuthnMakeCredential(makeCredentialBody);

    console.log({ creds });
  };

  return {
    register,
    onSubmit: handleSubmit(onSubmit),
  };
}

function bufToBase64(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer);
  let str = "";

  for (const charCode of bytes) {
    str += String.fromCharCode(charCode);
  }

  const base64String = btoa(str);

  return base64String.replace(/\+/g, "-").replace(/\//g, "_").replace(/=/g, "");
}
