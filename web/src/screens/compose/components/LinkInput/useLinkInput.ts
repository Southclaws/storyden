import { useFormContext } from "react-hook-form";

import { FormShape } from "../ComposeForm/useComposeForm";

export function useLinkInput() {
  const ctx = useFormContext<FormShape>();

  return {
    register: ctx.register,
    fieldError: ctx.formState.errors.category,
  };
}
