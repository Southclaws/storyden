"use client";

import { ProfileReference } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { MediaAddIcon } from "@/components/ui/icons/Media";
import { Box } from "@/styled-system/jsx";

import { EditAvatarTrigger } from "../EditAvatar/EditAvatarModal";

export type Props = {
  profile: ProfileReference;
};

export function MemberAvatarEditable({ profile }: Props) {
  return (
    <Box position="absolute" w="full" h="full">
      <EditAvatarTrigger profile={profile} asChild>
        <Button
          type="button"
          aria-label="Change avatar"
          position="absolute"
          top="0"
          left="0"
          w="full"
          h="full"
          minW="0"
          p="0"
          borderRadius="full"
          variant="subtle"
          color="text.inverse"
        >
          <MediaAddIcon width="6" height="6" />
        </Button>
      </EditAvatarTrigger>
    </Box>
  );
}
