import { Text } from "@/components/ui/text";
import { LStack, WStack } from "@/styled-system/jsx";

import { useLibraryPagePermissions } from "./permissions";

export function EditingDraftWarning() {
  const { isAllowedToDirectEdit } = useLibraryPagePermissions();

  const label = isAllowedToDirectEdit
    ? "Draft edits will be visible once applied."
    : "Draft edits will be visible once approved.";

  return (
    <LStack
      borderWidth="thin"
      borderStyle="dashed"
      borderColor="visibility.draft.border"
      borderRadius="sm"
      bgColor="background.inset"
      p="2"
      gap="0"
    >
      <Text variant="supporting" color="text.default" fontWeight="semibold">
        Editing Draft
      </Text>
      <Text as="span" variant="supporting">
        {label}
      </Text>
    </LStack>
  );
}
