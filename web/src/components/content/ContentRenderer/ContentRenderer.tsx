"use client";

import { Suspense, lazy } from "react";

import { css, cx } from "@/styled-system/css";
import { LStack } from "@/styled-system/jsx";

import type { ContentComposerProps } from "../composer-props";

import "./styles.css";

import { useRenderedContent } from "./RenderedContentContext";

const ClientContent = lazy(() =>
  import("./ContentRendererBody").then((module) => ({
    default: module.ContentRendererBody,
  })),
);

export function ContentRenderer({
  initialValue,
  value,
  className,
}: ContentComposerProps) {
  const html = value ?? initialValue ?? "<p></p>";
  const rendered = useRenderedContent(html);

  return (
    <LStack
      containerType="inline-size"
      className={cx(
        "content-composer content-composer--readonly typography",
        className,
      )}
      position="relative"
      w="full"
      gap="1"
      minHeight="8"
    >
      <div className={css({ height: "full", width: "full" })}>
        <div
          className={cx(
            "content-composer__content tiptap ProseMirror content-renderer",
            css({ height: "full", width: "full" }),
          )}
        >
          {rendered !== undefined ? (
            rendered
          ) : !html ? null : html === "<p></p>" ? (
            <p>
              <br className="ProseMirror-trailingBreak" />
            </p>
          ) : (
            <Suspense fallback={null}>
              <ClientContent html={html} />
            </Suspense>
          )}
        </div>
      </div>
    </LStack>
  );
}
