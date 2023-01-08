import { FieldValues, SubmitHandler, useForm } from "react-hook-form";
import {
  webAuthnMakeCredential,
  webAuthnRequestCredential,
} from "src/api/openapi/auth";

import { startRegistration } from "@simplewebauthn/browser";
import { PublicKeyCredentialCreationOptionsJSON } from "@simplewebauthn/typescript-types";
import { useToast } from "@chakra-ui/react";
import { APIError } from "src/api/openapi/schemas";
import { useRouter } from "next/router";

export function useWebAuthn() {
  const router = useRouter();
  const toast = useToast();
  const { register, handleSubmit } = useForm({
    //
  });

  const onSubmit: SubmitHandler<FieldValues> = async ({
    username,
  }: FieldValues) => {
    const response = await webAuthnRequestCredential(username);

    // TODO: OpenAPI spec for WebAuthn requests and responses.
    const publicKey = response[
      "publicKey"
    ] as PublicKeyCredentialCreationOptionsJSON;

    const credential = await startRegistration({
      ...publicKey,
      excludeCredentials: [],
    });

    webAuthnMakeCredential(credential)
      .then((account) => {
        console.info("webauthn success", { accountID: account.id });
        router.push("/");
      })
      .catch((e: APIError) => {
        toast({
          title: "Problem",
          status: "error",
          description: e.message,
        });
      });
  };

  return {
    register,
    onSubmit: handleSubmit(onSubmit),
  };
}
