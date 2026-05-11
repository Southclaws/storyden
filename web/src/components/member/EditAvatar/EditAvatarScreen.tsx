import { IconEditor } from "@/components/site/IconEditor/IconEditor";
import { Unready } from "@/components/site/Unready";
import { InfoIcon } from "@/components/ui/icons/Info";
import { useI18n } from "@/i18n/provider";
import { HStack, VStack } from "@/styled-system/jsx";

import { Props, useEditAvatar } from "./useEditAvatar";

export function EditAvatarScreen(props: Props) {
  const { ready, error, initialValue, handleSave } = useEditAvatar(props);
  const { t } = useI18n();
  if (!ready) {
    return <Unready error={error} />;
  }

  return (
    <VStack maxW="prose">
      <IconEditor
        isAvatar
        initialValue={initialValue}
        onUpload={handleSave}
        showPreviews
      />
      <HStack color="fg.subtle">
        <InfoIcon width="4" />
        <p>{t("You can pinch or use a mouse wheel to zoom/crop.")}</p>
      </HStack>
    </VStack>
  );
}
