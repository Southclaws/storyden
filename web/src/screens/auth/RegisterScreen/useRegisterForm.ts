"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { startRegistration } from "@simplewebauthn/browser";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { useAccountGet } from "src/api/openapi/accounts";
import {
  authPasswordSignup,
  webAuthnMakeCredential,
  webAuthnRequestCredential,
} from "src/api/openapi/auth";
import { APIError } from "src/api/openapi/schemas";
import { deriveError } from "src/utils/error";

export type Props = {
  webauthn: boolean;
};

const KindSchema = z.enum(["password", "webauthn"]);
type Kind = z.infer<typeof KindSchema>;

const UsernameSchema = z
  .string()
  .min(1, "Please enter a username.")
  .max(30, "Maximum length is 30 characters.")
  .toLowerCase()
  .regex(
    /^[a-z0-9_-]+$/g,
    "Username can only contain latin letters, numbers, dashes and underscores.",
  );

const PasswordSchema = z
  .string()
  .min(10, "Password must be at least 10 characters.");

const FormSchema = z.object({
  identifier: UsernameSchema,
  token: z.string().optional(), // Validated properly during password submission
});
const FormPasswordSchema = z.object({
  identifier: UsernameSchema,
  token: PasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function useRegisterForm() {
  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const { push } = useRouter();
  const { mutate } = useAccountGet();

  function handler(kind: Kind) {
    return handleSubmit((payload) => {
      switch (kind) {
        case "password":
          return handlePassword(payload);
        case "webauthn":
          return handleWebauthn(payload);
      }
    });
  }

  async function handlePassword(payload: Form) {
    const parsed = FormPasswordSchema.safeParse(payload);
    if (!parsed.success) {
      return setError("root", parsed.error);
    }

    await authPasswordSignup(parsed.data)
      .then(() => {
        push("/");
        mutate();
      })
      .catch((e: APIError) => setError("root", { message: deriveError(e) }));
  }

  async function handleWebauthn(payload: Form) {
    try {
      const { publicKey } = await webAuthnRequestCredential(payload.identifier);

      const credential = await startRegistration({
        ...publicKey,
      });

      await webAuthnMakeCredential(credential);

      push("/");
      mutate();
    } catch (error) {
      setError("root", { message: deriveError(error) });
    }
  }

  return {
    form: {
      register,
      handlePassword: handler("password"),
      handleWebauthn: handler("webauthn"),
      errors,
    },
  };
}
