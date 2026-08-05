import { useFormContext } from "react-hook-form";

import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { Input } from "@/components/ui/input";

import { FormShape } from "../ComposeForm/useComposeForm";

export function LinkInput() {
  const form = useFormContext<FormShape>();

  return (
    <FormControl>
      <Input
        placeholder="Share a link with your post..."
        type="url"
        {...form.register("url")}
      />
      <FormErrorText>{form.formState.errors.url?.message}</FormErrorText>
    </FormControl>
  );
}
