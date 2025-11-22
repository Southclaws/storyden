import Markdown from "react-markdown";

import { EditIcon } from "@/components/ui/icons/Edit";
import { ShowIcon } from "@/components/ui/icons/ShowIcon";
import { Switch } from "@/components/ui/switch";
import { Box, LStack, styled } from "@/styled-system/jsx";

import { ComposerTools } from "../ComposerTools";
import { ContentDragOverlay } from "../ContentDragOverlay";
import { ContentComposerProps } from "../composer-props";

import { useContentComposerMarkdown } from "./useContentComposerMarkdown";

export function ContentComposerMarkdown(props: ContentComposerProps) {
  const {
    value,
    previewHTML,
    showPreview,
    isDragging,
    isDragError,
    uploadingCount,
    textareaRef,
    getDragOverlayMessage,
    handleBufferChange,
    handleTogglePreview,
    handlePaste,
    handleDrop,
    handleDragOver,
    handleDragEnter,
    handleDragLeave,
  } = useContentComposerMarkdown(props);

  if (props.disabled) {
    return (
      <LStack position="relative" minHeight="8" maxHeight="fit">
        <Markdown className="typography">{value}</Markdown>
      </LStack>
    );
  }

  return (
    <LStack position="relative" minHeight="8" maxHeight="fit">
      <ComposerTools
        enabled={!props.disabled}
        icon={<ShowIcon />}
        expandedIcon={<EditIcon />}
        onClick={handleTogglePreview}
        workingCount={uploadingCount}
      >
        <Switch size="sm" checked={showPreview} onClick={handleTogglePreview}>
          Preview
        </Switch>
      </ComposerTools>

      {showPreview ? (
        <>
          {previewHTML ? (
            <styled.div
              className="typography"
              dangerouslySetInnerHTML={{ __html: previewHTML }}
            />
          ) : (
            <styled.p height="14" color="fg.muted" fontStyle="italic">
              empty...
            </styled.p>
          )}
        </>
      ) : (
        <>
          <styled.textarea
            ref={textareaRef}
            onChange={handleBufferChange}
            onPaste={handlePaste}
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragEnter={handleDragEnter}
            onDragLeave={handleDragLeave}
            value={value}
            lineHeight="relaxed"
            w="full"
            minHeight="0"
            resize="none"
            appearance="none"
            border="none"
            outline="none"
            color="fg.default"
            fontSize="md"
            transitionDuration="normal"
            transitionTimingFunction="default"
            _placeholder={{
              color: "fg.default",
            }}
            style={{
              border: "none",
              transitionProperty: "border-color, border-width",
              overflow: "hidden",
            }}
            placeholder="Write your heart out..."
          />
          {isDragging && (
            <ContentDragOverlay
              isError={isDragError}
              message={getDragOverlayMessage()}
            />
          )}
        </>
      )}
    </LStack>
  );
}
