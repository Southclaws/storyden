"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { authPasswordSignup } from "@/api/openapi-client/auth";
import { APIError } from "@/api/openapi-schema";
import { PasswordSchema, UsernameSchema } from "@/lib/auth/schemas";
import { deriveError } from "@/utils/error";

export type Props = {
  invitationID?: string;
};

const FormSchema = z.object({
  identifier: UsernameSchema,
  token: z.string().optional(), // Validated properly during password submission
});
const FormPasswordSchema = z.object({
  identifier: UsernameSchema,
  token: PasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function useRegisterHandleForm({ invitationID }: Props) {
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

    await authPasswordSignup(parsed.data, { invitation_id: invitationID })
      .then(() => {
        push("/");
        mutate();
      })
      .catch((e: APIError) => setError("root", { message: deriveError(e) }));
  }

  return {
    form: {
      register,
      handlePassword: handleSubmit(handlePassword),
      errors,
    },
  };
}
