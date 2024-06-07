import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { Input } from "@/components/ui/input";

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
