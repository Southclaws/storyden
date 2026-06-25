"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { authPasswordSignin } from "@/api/openapi-client/auth";
import { APIError } from "@/api/openapi-schema";
import { passkeyLogin } from "@/components/auth/webauthn/utils";
import { deriveError } from "@/utils/error";

import { ExistingPasswordSchema, UsernameSchema } from "@/lib/auth/schemas";
import { isWebauthnAvailable } from "@/lib/auth/webauthn";

export type Props = {
  webauthn: boolean;
};

const KindSchema = z.enum(["password", "webauthn"]);
type Kind = z.infer<typeof KindSchema>;

const FormSchema = z.object({
  identifier: UsernameSchema,
  token: z.string().optional(), // Validated properly during password submission
});
const FormPasswordSchema = z.object({
  identifier: UsernameSchema,
  token: ExistingPasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function useLoginHandleForm() {
  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const { push } = useRouter();
  const searchParams = useSearchParams();
  const returnURL = searchParams.get("return_url") ?? "/";
  const { mutate } = useAccountGet();

  const isWebauthnEnabled = isWebauthnAvailable();

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
      if (parsed.error.formErrors.fieldErrors.identifier) {
        setError("identifier", {
          message: parsed.error.formErrors.fieldErrors.identifier?.join(", "),
        });
      }

      if (parsed.error.formErrors.fieldErrors.token) {
        setError("token", {
          message: parsed.error.formErrors.fieldErrors.token?.join(", "),
        });
      }

      return;
    }

    await authPasswordSignin(parsed.data)
      .then(() => {
        push(returnURL);
        mutate();
      })
      .catch((e: APIError) => setError("root", { message: deriveError(e) }));
  }

  async function handleWebauthn(payload: Form) {
    try {
      await passkeyLogin(payload.identifier);
      push(returnURL);
      mutate();
    } catch (error) {
      setError("root", { message: deriveError(error) });
    }
  }

  return {
    form: {
      register,
      isWebauthnEnabled,
      handlePassword: handler("password"),
      handleWebauthn: handler("webauthn"),
      errors,
    },
  };
}
