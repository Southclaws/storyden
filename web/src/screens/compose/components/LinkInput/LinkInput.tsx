import { FormControl } from "src/theme/components/FormControl";
import { FormErrorText } from "src/theme/components/FormErrorText";
import { Input } from "src/theme/components/Input";

import { useLinkInput } from "./useLinkInput";

export function LinkInput() {
  const { register, fieldError } = useLinkInput();

  return (
    <FormControl>
      <Input
        size="xs"
        placeholder="Share a link with your post..."
        type="url"
        {...register("url")}
      />
      <FormErrorText>{fieldError?.message}</FormErrorText>
    </FormControl>
  );
}
