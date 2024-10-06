import { InfoIcon } from "lucide-react";

import { IconEditor } from "@/components/site/IconEditor/IconEditor";
import { Unready } from "@/components/site/Unready";
import { HStack, VStack } from "@/styled-system/jsx";

import { Props, useEditAvatar } from "./useEditAvatar";

export function EditAvatarScreen(props: Props) {
  const { ready, error, initialValue, handleSave } = useEditAvatar(props);
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
        <InfoIcon width="1rem" />
        <p>You can pinch or use a mouse wheel to zoom/crop.</p>
      </HStack>
    </VStack>
  );
}
