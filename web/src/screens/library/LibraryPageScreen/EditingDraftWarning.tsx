import { Heading } from "@/components/ui/heading";
import { LStack, WStack, styled } from "@/styled-system/jsx";

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
      <Heading size="sm">Editing Draft</Heading>
      <styled.span color="text.subtle" fontSize="sm">
        {label}
      </styled.span>
    </LStack>
  );
}
