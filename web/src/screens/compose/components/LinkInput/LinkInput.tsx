import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { Input } from "@/components/ui/input";
import { useI18n } from "@/i18n/provider";

import { useLinkInput } from "./useLinkInput";

export function LinkInput() {
  const { register, fieldError } = useLinkInput();
  const { t } = useI18n();

  return (
    <FormControl>
      <Input
        size="xs"
        placeholder={t("Share a link with your post...")}
        type="url"
        {...register("url")}
      />
      <FormErrorText>{fieldError?.message}</FormErrorText>
    </FormControl>
  );
}
