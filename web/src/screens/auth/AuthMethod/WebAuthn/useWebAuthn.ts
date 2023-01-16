import { FieldValues, SubmitHandler, useForm } from "react-hook-form";
import {
  webAuthnGetAssertion,
  webAuthnMakeAssertion,
  webAuthnMakeCredential,
  webAuthnRequestCredential,
} from "src/api/openapi/auth";

import {
  startAuthentication,
  startRegistration,
} from "@simplewebauthn/browser";
import { PublicKeyCredentialCreationOptionsJSON } from "@simplewebauthn/typescript-types";
import { useToast } from "@chakra-ui/react";
import { APIError } from "src/api/openapi/schemas";
import { useRouter } from "next/router";
import * as z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { errorToast } from "src/components/ErrorBanner";

export const FormSchema = z.object({
  username: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function useWebAuthn() {
  const router = useRouter();
  const toast = useToast();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  async function signin({ username }: Form) {
    try {
      const response = await webAuthnGetAssertion(username);

      // TODO: OpenAPI spec for WebAuthn requests and responses.
      const publicKey = response[
        "publicKey"
      ] as PublicKeyCredentialCreationOptionsJSON;

      const credential = await startAuthentication({
        ...publicKey,
      });

      // HACK:
      // 1. https://github.com/MasterKale/SimpleWebAuthn/issues/330
      // 2. https://github.com/go-webauthn/webauthn/issues/93
      // credential.response.userHandle =
      //   credential.response.userHandle?.replaceAll("=", "");

      console.log({ username, publicKey, credential });

      await webAuthnMakeAssertion(credential);

      router.push("/");
    } catch (error) {
      errorToast(toast)(error as APIError);
    }
  }

  async function signup({ username }: Form) {
    try {
      const response = await webAuthnRequestCredential(username);

      // TODO: OpenAPI spec for WebAuthn requests and responses.
      const publicKey = response[
        "publicKey"
      ] as PublicKeyCredentialCreationOptionsJSON;

      const credential = await startRegistration({
        ...publicKey,
        excludeCredentials: [],
      });

      await webAuthnMakeCredential(credential);

      router.push("/");
    } catch (error) {
      errorToast(toast)(error as APIError);
    }
  }

  function onSubmit(action: "signin" | "signup") {
    return action === "signin" ? handleSubmit(signin) : handleSubmit(signup);
  }

  return {
    register,
    onSubmit,
    errors,
  };
}
