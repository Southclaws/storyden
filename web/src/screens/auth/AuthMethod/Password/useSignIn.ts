import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { api } from "src/api";
import { Form, FormSchema } from "./common";

export default function useSignIn() {
  const { register, handleSubmit } = useForm({
    resolver: zodResolver(FormSchema),
  });

  const onSubmit = (payload: Form) => {
    api("/v1/auth/password/signin");
  };

  return {
    form: {
      register,
      handleSubmit: handleSubmit(onSubmit),
    },
  };
}
