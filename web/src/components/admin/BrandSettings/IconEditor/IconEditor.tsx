import AvatarEditor from "react-avatar-editor";

import { Button } from "@/components/ui/button";
import { MediaAddIcon } from "@/components/ui/icons/Media";
import { SaveIcon } from "@/components/ui/icons/Save";
import { Box, HStack, VStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { Props, useIconEditor } from "./useIconEditor";

const editorStyle = { backgroundColor: "var(--colors-gray-100)" };

export function IconEditor(props: Props) {
  const {
    ref,
    position,
    setPosition,
    onFileChange,
    onImageChange,
    onSave,
    saving,
    file,
    preview,
  } = useIconEditor(props);

  return (
    <VStack alignItems="start" w="min">
      <>
        <styled.input
          id="file-input"
          display="none"
          width="min"
          type="file"
          bgColor="background.inset"
          borderRadius="md"
          border="none"
          onChange={onFileChange}
        />
        <HStack w="full">
          <styled.label
            htmlFor="file-input"
            w="full"
            className={button({ variant: "outline" })}
          >
            <MediaAddIcon /> Upload icon
          </styled.label>
          <Button variant="solid" w="full" onClick={onSave} disabled={saving}>
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
            onImageChange={onImageChange}
          />
        </Box>
        <VStack alignItems="start" gap="2" justifyContent="space-between">
          <HStack>
            <Box
              w="8"
              h="8"
              borderRadius="md"
              overflow="hidden"
              backgroundColor="background.inset"
            >
              {preview && <styled.img src={preview} alt="" w="full" h="full" />}
            </Box>

            <Box
              w="8"
              h="8"
              borderRadius="full"
              overflow="hidden"
              backgroundColor="background.inset"
            >
              {preview && <styled.img src={preview} alt="" w="full" h="full" />}
            </Box>
          </HStack>

          <Box
            w="24"
            h="24"
            borderRadius="lg"
            overflow="hidden"
            backgroundColor="background.inset"
          >
            {preview && <styled.img src={preview} alt="" w="full" h="full" />}
          </Box>
        </VStack>
      </HStack>
    </VStack>
  );
}
