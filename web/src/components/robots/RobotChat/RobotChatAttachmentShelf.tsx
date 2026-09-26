import { AssetThumbnail } from "@/components/asset/AssetThumbnail";
import { IconButton } from "@/components/ui/icon-button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { Spinner } from "@/components/ui/spinner";
import { Box, HStack, styled } from "@/styled-system/jsx";

import { RobotChatImageAttachment } from "./useRobotChatImageAttachments";

type Props = {
  attachments: RobotChatImageAttachment[];
  onRemove: (id: string) => void;
};

export function RobotChatAttachmentShelf({ attachments, onRemove }: Props) {
  if (attachments.length === 0) {
    return null;
  }

  const assets = attachments.flatMap((attachment) =>
    attachment.asset ? [attachment.asset] : [],
  );

  return (
    <HStack
      aria-label="Attached images"
      w="full"
      gap="2"
      overflowX="auto"
      overflowY="hidden"
      px="3"
      pt="3"
      pb="1"
    >
      {attachments.map((attachment) => (
        <Box
          key={attachment.id}
          position="relative"
          flexShrink="0"
          w="20"
          h="20"
          overflow="hidden"
          borderRadius="md"
          background="background.controlInset"
        >
          {attachment.asset ? (
            <AssetThumbnail
              asset={attachment.asset}
              set={assets}
              setIndex={assets.findIndex(
                (asset) => asset.id === attachment.asset?.id,
              )}
              showDeleteButton
              deleteLabel={`Remove ${attachment.filename}`}
              handleDelete={() => onRemove(attachment.id)}
            />
          ) : (
            <>
              <styled.img
                src={attachment.previewURL}
                alt=""
                w="full"
                h="full"
                objectFit="cover"
                opacity="6"
              />
              <Box
                role="status"
                aria-label={`Uploading ${attachment.filename}`}
                position="absolute"
                inset="0"
                display="flex"
                alignItems="center"
                justifyContent="center"
                color="text.default"
              >
                <Spinner />
              </Box>
              <IconButton
                type="button"
                position="absolute"
                top="1"
                right="1"
                intent="destructive"
                variant="subtle"
                w="5"
                h="5"
                minW="5"
                aria-label={`Remove ${attachment.filename}`}
                onClick={() => onRemove(attachment.id)}
              >
                <DeleteIcon />
              </IconButton>
            </>
          )}
        </Box>
      ))}
    </HStack>
  );
}
