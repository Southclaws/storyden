import { Button, HStack, VStack } from "@chakra-ui/react";
import AvatarEditor from "react-avatar-editor";

import { styled } from "@/styled-system/jsx";

import { Props, useIconEditor } from "./useIconEditor";

const editorStyle = { backgroundColor: "var(--chakra-colors-gray-100)" };

export function IconEditor(props: Props) {
  const { position, setPosition, onFileChange, onSave, file } =
    useIconEditor(props);

  return (
    <VStack align="start">
      <HStack justify="start" align="end">
        <AvatarEditor
          image={file}
          width={136}
          height={136}
          borderRadius={12}
          style={editorStyle}
          border={0}
          color={[255, 255, 255, 1]}
          scale={1}
          position={position}
          onPositionChange={setPosition}
        />
        <VStack align="start">
          <HStack>
            <AvatarEditor
              image={file}
              width={32}
              height={32}
              borderRadius={6}
              style={editorStyle}
              border={0}
              color={[255, 255, 255, 1]}
              scale={1}
              position={position}
            />
            <AvatarEditor
              image={file}
              width={32}
              height={32}
              borderRadius={99}
              style={editorStyle}
              border={0}
              color={[255, 255, 255, 1]}
              scale={1}
              position={position}
            />
          </HStack>

          <AvatarEditor
            image={file}
            width={96}
            height={96}
            borderRadius={12}
            style={editorStyle}
            border={0}
            color={[255, 255, 255, 1]}
            scale={1}
            position={position}
          />
        </VStack>
      </HStack>
      <>
        <styled.input
          id="file-input"
          display="none"
          width="min-content"
          type="file"
          bgColor="gray.100"
          borderRadius="md"
          border={0}
          onChange={onFileChange}
        />
        <HStack>
          <Button as="label" htmlFor="file-input" variant="outline">
            Edit icon
          </Button>
          <Button onClick={onSave}>Save icon</Button>
        </HStack>
      </>
    </VStack>
  );
}
