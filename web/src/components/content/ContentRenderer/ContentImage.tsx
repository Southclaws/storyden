"use client";

import type { ElementType } from "react";

import { Button } from "@/components/ui/button";
import { ProgressCircle } from "@/components/ui/progress";
import { Text } from "@/components/ui/text";
import { css } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";
import { getImageVariantURL } from "@/utils/asset";

export function ContentImage({
  attrs,
  isEditable = false,
  isSelected = false,
  onRetry,
  onCancel,
  Wrapper = "div",
}: {
  attrs: Record<string, any>;
  isEditable?: boolean;
  isSelected?: boolean;
  onRetry?: () => void;
  onCancel?: () => void;
  Wrapper?: ElementType;
}) {
  const isUploading = attrs["data-uploading"] === "true";
  const uploadError = attrs["data-upload-error"];
  const uploadProgress = attrs["data-upload-progress"];
  const progressPercent = uploadProgress ? parseInt(uploadProgress, 10) : 0;
  const dataDiff = attrs["data-diff"] as "insertion" | "deletion" | undefined;

  const diffClass = dataDiff ? `diff-node-${dataDiff}` : "";

  return (
    <Wrapper
      data-node-view-wrapper=""
      style={{ whiteSpace: "normal" }}
      data-diff={dataDiff}
      className={`${diffClass} ${css({
        position: "relative",
        display: "inline-block",
        cursor: "pointer",
        outlineWidth: isSelected ? "medium" : "none",
        outlineStyle: "solid",
        outlineColor: isSelected ? "blue.a6" : "transparent",
        borderRadius: "lg",
        userSelect: isEditable ? "none" : "auto",
        saturate: isSelected && !isUploading ? "150%" : "100%",
        filter: "auto",
        background: isSelected && !isUploading ? "blue.5" : "transparent",
        mixBlendMode: isSelected && !isUploading ? "screen" : "normal",
      })}`}
    >
      <styled.img
        borderRadius="lg"
        opacity={isUploading ? "5" : "full"}
        transition="all"
        {...attrs}
        src={getImageVariantURL(attrs["src"], "l")}
      />
      {isUploading && (
        <styled.div
          position="absolute"
          top="0"
          left="0"
          width="full"
          height="full"
          display="flex"
          flexDirection="column"
          alignItems="center"
          justifyContent="center"
          pointerEvents="none"
          gap="3"
          padding="4"
        >
          <ProgressCircle value={progressPercent} size="md" />
        </styled.div>
      )}
      {uploadError && (
        <styled.div
          position="absolute"
          inset="0"
          display="flex"
          flexDirection="column"
          alignItems="center"
          justifyContent="center"
          backgroundColor="status.danger.surface"
          opacity="9"
          borderRadius="md"
          padding="3"
          gap="2"
          userSelect="none"
          contentEditable={false}
        >
          <Text
            variant="supporting"
            color="status.danger.content"
            fontWeight="medium"
          >
            Upload failed
          </Text>
          <styled.div display="flex" gap="2">
            <Button type="button" variant="outline" onClick={onRetry}>
              Retry
            </Button>
            <Button type="button" variant="ghost" onClick={onCancel}>
              Remove
            </Button>
          </styled.div>
        </styled.div>
      )}
    </Wrapper>
  );
}
