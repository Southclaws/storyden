import { CSSProperties } from "react";
import AvatarEditor from "react-avatar-editor";

import { Button } from "@/components/ui/button";
import { Flex, LStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { Props, useAssetUploadEditor } from "./useAssetUploadEditor";

export interface AssetUploadEditorProps extends Props {
  width?: number;
  height?: number;
}

const DEFAULT_WIDTH = 1280;
const DEFAULT_HEIGHT = 720;

export function AssetUploadEditor(props: AssetUploadEditorProps) {
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
  } = useAssetUploadEditor(props);

  const width = props.width ?? DEFAULT_WIDTH;
  const height = props.height ?? DEFAULT_HEIGHT;
  const aspectRatio = `${width} / ${height}`;

  return (
    <LStack
      id="asset-upload-editor"
      width="full"
      borderColor="border.default"
      borderWidth="thin"
      borderRadius="xl"
    >
      <LStack overflow="hidden" w="full" gap="0">
        <styled.input
          id="asset-upload-editor__file-input"
          display="none"
          type="file"
          accept="image/*"
          onChange={onFileChange}
        />

        <styled.div
          ref={pinchCaptureRef}
          width="full"
          aspectRatio={aspectRatio}
          backgroundColor="bg.muted"
          borderTopRadius="xl"
          display="flex"
          alignItems="center"
          justifyContent="center"
          onPointerDown={(e) => e.stopPropagation()}
          onTouchStart={(e) => e.stopPropagation()}
        >
          <AvatarEditor
            ref={ref}
            image={file}
            width={width}
            height={height}
            style={
              {
                width: "100%",
                height: "auto",
                margin: 0,
                backgroundColor: "transparent",
              } as CSSProperties
            }
            border={0}
            color={[255, 255, 255, 1]}
            scale={scale}
            position={position}
            onPositionChange={onPositionChange}
            crossOrigin="anonymous"
          />
        </styled.div>

        <Flex flexDirection="row" w="full" marginTop="0" gap="0">
          <styled.label
            flexGrow="1"
            borderRightRadius="none"
            borderTopRadius="none"
            borderRadius="xl"
            className={button({
              variant: "ghost",
            })}
            htmlFor="asset-upload-editor__file-input"
          >
            Edit
          </styled.label>
          <Button
            flexGrow="1"
            borderLeftRadius="none"
            borderTopRadius="none"
            borderRadius="xl"
            onClick={onSave}
            disabled={saving}
          >
            Save
          </Button>
        </Flex>
      </LStack>
    </LStack>
  );
}

// Legacy export for backward compatibility
export const CategoryImageEditor = AssetUploadEditor;
