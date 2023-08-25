import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  startAuthentication,
  startRegistration,
} from "@simplewebauthn/browser";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { useAccountGet } from "src/api/openapi/accounts";
import {
  webAuthnGetAssertion,
  webAuthnMakeAssertion,
  webAuthnMakeCredential,
  webAuthnRequestCredential,
} from "src/api/openapi/auth";
import { APIError } from "src/api/openapi/schemas";
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
  const { mutate } = useAccountGet();

  async function signin({ username }: Form) {
    try {
      const { publicKey } = await webAuthnGetAssertion(username);
      const credential = await startAuthentication(publicKey);

      // HACK:
      // 1. https://github.com/MasterKale/SimpleWebAuthn/issues/330
      // 2. https://github.com/go-webauthn/webauthn/issues/93
      credential.response.userHandle = undefined;

      await webAuthnMakeAssertion(credential);

      router.push("/");
      mutate();
    } catch (error) {
      errorToast(toast)(error as APIError);
    }
  }

  async function signup({ username }: Form) {
    try {
      const { publicKey } = await webAuthnRequestCredential(username);

      const credential = await startRegistration({
        ...publicKey,
      });

      await webAuthnMakeCredential(credential);

      router.push("/");
      mutate();
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
