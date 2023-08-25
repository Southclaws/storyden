import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { useAccountGet } from "src/api/openapi/accounts";
import { authPasswordSignin, authPasswordSignup } from "src/api/openapi/auth";
import { APIError } from "src/api/openapi/schemas";

export const FormSchema = z.object({
  identifier: z.string(),
  token: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function usePassword() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const toast = useToast();
  const { push } = useRouter();
  const { mutate } = useAccountGet();

  async function signin(payload: Form) {
    await authPasswordSignin(payload)
      .then(() => {
        push("/");
        mutate();
      })
      .catch((e: APIError) =>
        toast({
          title: "Failed!",
          description: `Sign in failed: ${e.message}`,
          status: "error",
        })
      );
  }

  async function signup(payload: Form) {
    await authPasswordSignup(payload)
      .then(() => {
        push("/");
        mutate();
      })
      .catch((e: APIError) =>
        toast({
          title: "Failed!",
          description: `Sign up failed: ${e.message}`,
          status: "error",
        })
      );
  }

  function onSubmit(action: "signin" | "signup") {
    return action === "signin" ? handleSubmit(signin) : handleSubmit(signup);
  }

  return {
    form: {
      register,
      onSubmit: onSubmit,
      errors,
    },
  };
}
