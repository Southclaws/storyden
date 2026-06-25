import Link from "next/link";
import { ReactNode } from "react";
import Markdown, { Components } from "react-markdown";

import { CalendarIcon } from "@/components/ui/icons/Calendar";
import { CardIcon } from "@/components/ui/icons/Card";
import { CollectionIcon } from "@/components/ui/icons/Collection";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { EditIcon } from "@/components/ui/icons/Edit";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { ProfileIcon } from "@/components/ui/icons/Profile";
import { ReplyIcon } from "@/components/ui/icons/Reply";
import { ShowIcon } from "@/components/ui/icons/ShowIcon";
import { Switch } from "@/components/ui/switch";
import { css } from "@/styled-system/css";
import { LStack, styled } from "@/styled-system/jsx";
import { markdownURLTransform, remarkLooseLists } from "@/utils/markdown";

import { ComposerTools } from "../ComposerTools";
import { ContentDragOverlay } from "../ContentDragOverlay";
import { ContentComposerProps } from "../composer-props";

import { useContentComposerMarkdown } from "./useContentComposerMarkdown";

const markdownComponents: Components = {
  a: ({ href, children }) => {
    if (!href) {
      return <a>{children}</a>;
    }

    const ref = parseSDRHref(href);
    if (!ref) {
      return <a href={href}>{children}</a>;
    }

    return (
      <SDRInlineReference kind={ref.kind} id={ref.id}>
        {children}
      </SDRInlineReference>
    );
  },
};

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
      <LStack
        className="markdown-editor-readonly"
        position="relative"
        minHeight="8"
        maxHeight="fit"
      >
        <Markdown
          className="typography"
          components={markdownComponents}
          remarkPlugins={[remarkLooseLists]}
          urlTransform={markdownURLTransform}
        >
          {value}
        </Markdown>
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

function SDRInlineReference({
  kind,
  id,
  children,
}: {
  kind: string;
  id: string;
  children: ReactNode;
}) {
  return (
    <Link
      href={`/_/resolve/${kind}/${id}`}
      className={css({
        display: "inline-flex",
        alignItems: "center",
        gap: "1",
        px: "2",
        verticalAlign: "baseline",
        borderRadius: "full",
        backgroundColor: "bg.muted",
        color: "fg.default",
        textDecoration: "none",
        fontWeight: "medium",
        fontSize: "xs",
        whiteSpace: "nowrap",
        _hover: {
          backgroundColor: "bg.subtle",
          textDecoration: "none",
        },
      })}
    >
      {renderSDRIcon(kind)}
      <span>{children}</span>
    </Link>
  );
}

function parseSDRHref(href: string): { kind: string; id: string } | null {
  const match = /^sdr:([a-z]+)\/([a-z0-9]+)$/i.exec(href);
  if (!match) {
    return null;
  }

  return {
    kind: match[1]!,
    id: match[2]!,
  };
}

function renderSDRIcon(kind: string) {
  switch (kind) {
    case "node":
      return <LibraryIcon width="3" height="3" aria-hidden />;
    case "thread":
    case "post":
      return <DiscussionIcon width="3" height="3" aria-hidden />;
    case "reply":
      return <ReplyIcon width="3" height="3" aria-hidden />;
    case "profile":
      return <ProfileIcon width="3" height="3" aria-hidden />;
    case "collection":
      return <CollectionIcon width="3" height="3" aria-hidden />;
    case "event":
      return <CalendarIcon width="3" height="3" aria-hidden />;
    default:
      return <CardIcon width="3" height="3" aria-hidden />;
  }
}
