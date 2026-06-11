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
      bgColor="bg.subtle"
      p="2"
      gap="0"
    >
      <Heading size="sm">Editing Draft</Heading>
      <styled.span color="fg.muted" fontSize="sm">
        {label}
      </styled.span>
    </LStack>
  );
}
