import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/router";
import { useForm } from "react-hook-form";
import { authPasswordSignin, authPasswordSignup } from "src/api/openapi/auth";
import { APIError } from "src/api/openapi/schemas";
import { Form, FormSchema } from "./common";

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

  async function signin(payload: Form, x: any) {
    await authPasswordSignin(payload)
      .then(() => {
        push("/");
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
