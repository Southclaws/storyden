import AvatarEditor from "react-avatar-editor";

import { Button } from "@/components/ui/button";
import { MediaAddIcon } from "@/components/ui/icons/Media";
import { SaveIcon } from "@/components/ui/icons/Save";
import { Box, HStack, VStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { Props, useIconEditor } from "./useIconEditor";

const editorStyle = { backgroundColor: "var(--colors-gray-100)" };

export function IconEditor(props: Props) {
  const { ref, position, setPosition, onFileChange, onSave, saving, file } =
    useIconEditor(props);

  return (
    <VStack alignItems="start" w="min">
      <>
        <styled.input
          id="file-input"
          display="none"
          width="min"
          type="file"
          bgColor="bg.subtle"
          borderRadius="md"
          border="none"
          onChange={onFileChange}
        />
        <HStack w="full">
          <styled.label
            htmlFor="file-input"
            w="full"
            className={button({ size: "xs", variant: "outline" })}
          >
            <MediaAddIcon /> Upload icon
          </styled.label>
          <Button
            size="xs"
            variant="solid"
            w="full"
            onClick={onSave}
            disabled={saving}
          >
            <SaveIcon /> Save icon
          </Button>
        </HStack>
      </>

      <HStack>
        <Box borderRadius="lg" overflow="hidden">
          <AvatarEditor
            ref={ref}
            image={file}
            width={136}
            height={136}
            style={editorStyle}
            border={0}
            color={[255, 255, 255, 1]}
            scale={1}
            position={position}
            onPositionChange={saving ? undefined : setPosition}
          />
        </Box>
        <VStack alignItems="start" gap="2" justifyContent="space-between">
          <HStack>
            <Box borderRadius="md" overflow="hidden">
              <AvatarEditor
                image={file}
                width={32}
                height={32}
                style={editorStyle}
                border={0}
                color={[255, 255, 255, 1]}
                scale={1}
                position={position}
              />
            </Box>

            <Box borderRadius="full" overflow="hidden">
              <AvatarEditor
                image={file}
                width={32}
                height={32}
                style={editorStyle}
                border={0}
                color={[255, 255, 255, 1]}
                scale={1}
                position={position}
              />
            </Box>
          </HStack>

          <Box borderRadius="lg" overflow="hidden">
            <AvatarEditor
              image={file}
              width={96}
              height={96}
              style={editorStyle}
              border={0}
              color={[255, 255, 255, 1]}
              scale={1}
              position={position}
            />
          </Box>
        </VStack>
      </HStack>
    </VStack>
  );
}
