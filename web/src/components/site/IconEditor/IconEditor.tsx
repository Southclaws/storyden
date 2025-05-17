import { CSSProperties } from "react";
import AvatarEditor from "react-avatar-editor";

import { avatarSize } from "@/components/member/MemberBadge/MemberAvatar";
import { Button } from "@/components/ui/button";
import { ArrowLeftIcon } from "@/components/ui/icons/Arrow";
import { Box, Flex, HStack, LStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { Props, useIconEditor } from "./useIconEditor";

export function IconEditor(props: Props) {
  const {
    ref,
    pinchCaptureRef,
    scale,
    position,
    onPositionChange,
    onFileChange,
    onSave,
    saving,
    file,
    isDirty,
  } = useIconEditor(props);

  const { showPreviews } = props;

  return (
    <LStack
      id="icon-editor"
      width="full"
      borderColor="border.default"
      borderWidth="thin"
      borderRadius="xl"
    >
      {showPreviews && (
        <HStack pl="2" pt="2" color="fg.subtle">
          <AvatarEditor
            image={file}
            style={{
              borderRadius: "100%",
              ...avatarSize("sm"),
            }}
            border={0}
            color={[255, 255, 255, 1]}
            scale={scale}
            position={position}
            onPositionChange={onPositionChange}
          />
          <ArrowLeftIcon width="4" />
          <p>How it&apos;ll look on posts</p>
        </HStack>
      )}

      <LStack overflow="hidden" w="full" gap="0">
        <styled.input
          id="icon-editor__file-input"
          display="none"
          type="file"
          onChange={onFileChange}
        />

        <Box w="full" ref={pinchCaptureRef}>
          <AvatarEditor
            ref={ref}
            image={file}
            style={
              {
                width: "100%",
                height: "unset",
                aspectRatio: "1.0",
                margin: 0,
                backgroundColor: "var(--colors-gray-2)",
              } as CSSProperties
            }
            border={0}
            color={[255, 255, 255, 1]}
            scale={scale}
            position={position}
            onPositionChange={onPositionChange}
          />
        </Box>

        <Flex flexDirection="row" w="full" marginTop="0" gap="0">
          <styled.label
            flexGrow="1"
            borderRightRadius="none"
            borderTopRadius="none"
            borderRadius="xl"
            className={button({
              variant: "ghost",
            })}
            htmlFor="icon-editor__file-input"
          >
            Edit
          </styled.label>
          <Button
            flexGrow="1"
            borderLeftRadius="none"
            borderTopRadius="none"
            borderRadius="xl"
            onClick={onSave}
            disabled={!isDirty || saving}
          >
            Save
          </Button>
        </Flex>
      </LStack>
    </LStack>
  );
}
