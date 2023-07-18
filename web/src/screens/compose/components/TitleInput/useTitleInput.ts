import { useFormContext } from "react-hook-form";

import { FormShape } from "../ComposeForm/useComposeForm";

export function useTitleInput() {
  const ctx = useFormContext<FormShape>();

  return {
    control: ctx.control,
    fieldError: ctx.formState.errors["title"],
  };
}
