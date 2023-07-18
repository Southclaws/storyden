import { useFormContext } from "react-hook-form";

import { FormShape } from "../ComposeForm/useComposeForm";

export function useBodyInput() {
  const ctx = useFormContext<FormShape>();

  return {
    control: ctx.control,
    fieldError: ctx.formState.errors.category,
  };
}
