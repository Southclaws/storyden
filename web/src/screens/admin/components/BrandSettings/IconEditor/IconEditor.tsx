import AvatarEditor from "react-avatar-editor";

import { Button } from "src/theme/components/Button";

import { Box, HStack, VStack, styled } from "@/styled-system/jsx";

import { Props, useIconEditor } from "./useIconEditor";

const editorStyle = { backgroundColor: "var(--colors-gray-100)" };

export function IconEditor(props: Props) {
  const { ref, position, setPosition, onFileChange, onSave, saving, file } =
    useIconEditor(props);

  return (
    <VStack alignItems="start">
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
      <>
        <styled.input
          id="file-input"
          display="none"
          width="min"
          type="file"
          bgColor="gray.100"
          borderRadius="md"
          border="none"
          onChange={onFileChange}
        />
        <HStack>
          <styled.label htmlFor="file-input">Edit icon</styled.label>
          <Button onClick={onSave} disabled={saving}>
            Save icon
          </Button>
        </HStack>
      </>
    </VStack>
  );
}
