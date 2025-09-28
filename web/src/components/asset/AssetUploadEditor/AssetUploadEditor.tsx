import { CSSProperties } from "react";
import AvatarEditor from "react-avatar-editor";

import { Button } from "@/components/ui/button";
import { Flex, LStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { Props, useAssetUploadEditor } from "./useAssetUploadEditor";

export interface AssetUploadEditorProps extends Props {
  aspectRatio?: string;
  width?: number;
  height?: number;
}

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

  const aspectRatio = props.aspectRatio ?? "16 / 9";
  const width = props.width ?? 400;
  const height = props.height ?? Math.round(width / (16 / 9));

  return (
    <LStack
      id="asset-upload-editor"
      width="full"
      style={
        {
          "--max-width": `${width}px`,
        } as any
      }
      // maxWidth="var(--max-width)"
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
        >
          <AvatarEditor
            ref={ref}
            image={file}
            width={width}
            height={height}
            style={
              {
                maxWidth: "100%",
                maxHeight: "100%",
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
