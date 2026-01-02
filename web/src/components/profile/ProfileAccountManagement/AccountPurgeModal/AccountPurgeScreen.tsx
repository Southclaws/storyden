import { ModerationActionPurgeAccountContentType } from "@/api/openapi-schema";
import { FormControl } from "@/components/ui/FormControl";
import { FormLabel } from "@/components/ui/FormLabel";
import { Button } from "@/components/ui/button";
import { CardGroupSelect } from "@/components/ui/form/CardGroupSelect";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Props, useAccountPurgeScreen } from "./useAccountPurge";

const CONTENT_TYPE_LABELS: Record<
  ModerationActionPurgeAccountContentType,
  { name: string; description: string }
> = {
  [ModerationActionPurgeAccountContentType.threads]: {
    name: "Threads",
    description: "Delete all threads created by this account",
  },
  [ModerationActionPurgeAccountContentType.replies]: {
    name: "Replies",
    description: "Delete all replies made by this account",
  },
  [ModerationActionPurgeAccountContentType.reacts]: {
    name: "Reactions",
    description: "Remove all reactions made by this account",
  },
  [ModerationActionPurgeAccountContentType.likes]: {
    name: "Likes",
    description: "Remove all likes made by this account",
  },
  [ModerationActionPurgeAccountContentType.nodes]: {
    name: "Library Pages",
    description: "Delete all library pages created by this account",
  },
  [ModerationActionPurgeAccountContentType.collections]: {
    name: "Collections",
    description: "Delete all collections created by this account",
  },
  [ModerationActionPurgeAccountContentType.profile_bio]: {
    name: "Profile Bio",
    description: "Clear the account's profile biography",
  },
};

const CONTENT_TYPES = Object.entries(
  ModerationActionPurgeAccountContentType,
).map(([_, value]) => ({
  value,
  label: CONTENT_TYPE_LABELS[value].name,
  description: CONTENT_TYPE_LABELS[value].description,
}));

export function AccountPurgeScreen(props: Props) {
  const {
    form,
    handlers: { handlePurge },
  } = useAccountPurgeScreen(props);

  return (
    <styled.form
      className={lstack()}
      h="full"
      justifyContent="space-between"
      onSubmit={handlePurge}
    >
      <LStack px="0.5" maxH="full" pb="1" overflowY="scroll">
        <Box
          p="3"
          borderRadius="sm"
          borderWidth="thin"
          borderColor="border.warning"
          bgColor="bg.warning"
        >
          <HStack gap="2" color="fg.warning">
            <WarningIcon w="5" flexShrink="0" />
            <LStack gap="1">
              <styled.p fontWeight="semibold" fontSize="sm">
                Destructive Action
              </styled.p>
              <styled.p fontSize="xs">
                This will permanently delete the selected content types from
                this account. This action cannot be undone.
              </styled.p>
            </LStack>
          </HStack>
        </Box>

        <FormControl>
          <FormLabel>Content Types to Purge</FormLabel>
          <CardGroupSelect
            control={form.control}
            name="contentTypes"
            items={CONTENT_TYPES}
          />
          <styled.p fontSize="xs" color="fg.subtle" mt="1">
            Select the types of content you want to purge from this account
          </styled.p>
        </FormControl>
      </LStack>

      <WStack>
        <Button
          flexGrow="1"
          variant="solid"
          disabled={!form.formState.isDirty || form.formState.isSubmitting}
          loading={form.formState.isSubmitting}
          type="submit"
        >
          Purge Selected Content
        </Button>
      </WStack>
    </styled.form>
  );
}
