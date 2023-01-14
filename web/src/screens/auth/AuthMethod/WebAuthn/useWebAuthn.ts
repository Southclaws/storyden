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
    const response = await webAuthnGetAssertion(username);

    // TODO: OpenAPI spec for WebAuthn requests and responses.
    const publicKey = response[
      "publicKey"
    ] as PublicKeyCredentialCreationOptionsJSON;

    console.log({ response, publicKey });

    const credential = await startAuthentication({
      ...publicKey,
    });

    console.log({ credential });

    const r = await webAuthnMakeAssertion(credential);
    console.log({ r });

    // router.push("/");
  }

  async function signup({ username }: Form) {
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
  }

  function onSubmit(action: "signin" | "signup") {
    try {
      return action === "signin" ? handleSubmit(signin) : handleSubmit(signup);
    } catch (error) {
      errorToast(toast)(error as APIError);
    }
    return null;
  }

  return {
    register,
    onSubmit,
    errors,
  };
}
