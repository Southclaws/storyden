import { FormControl, FormErrorMessage } from "src/theme/components";
import { Input } from "src/theme/components/Input";

import { useLinkInput } from "./useLinkInput";

export function LinkInput() {
  const { register, fieldError } = useLinkInput();

  return (
    <FormControl isInvalid={!!fieldError}>
      <Input
        size="xs"
        placeholder="Share a link with your post..."
        type="url"
        {...register("url")}
      />
      <FormErrorMessage>{fieldError?.message}</FormErrorMessage>
    </FormControl>
  );
}
