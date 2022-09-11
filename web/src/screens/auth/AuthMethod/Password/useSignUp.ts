import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/router";
import { useForm } from "react-hook-form";
import { authPasswordSignup } from "src/api/openapi/auth";

import { Form, FormSchema } from "./common";

export default function useSignUp() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  const toast = useToast();

  const { push } = useRouter();

  const onSubmit = async (payload: Form) => {
    await authPasswordSignup(payload)
      .then(() => {
        push("/");
      })
      .catch((e) =>
        toast({
          title: "Failed!",
          description: `Failed to sign up: ${e}`,
          status: "error",
        })
      );
  };

  return {
    form: {
      register,
      onSubmit: handleSubmit(onSubmit),
      errors,
    },
  };
}
