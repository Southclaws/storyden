"use client";

import dynamic from "next/dynamic";
import type { PropsWithChildren } from "react";

import { contentOutline } from "@/styled-system/recipes";

import type { CollectOutlineEntries } from "./outline";
import { useContentOutline } from "./useContentOutline";

const ContentOutlineRail = dynamic(
  () =>
    import("./ContentOutline.internal").then((mod) => mod.ContentOutlineRail),
  { ssr: false },
);

export type ContentOutlineProps = PropsWithChildren<{
  label: string;
  collectEntries: CollectOutlineEntries;
  disabled?: boolean;
  placement?: "inline" | "gutter";
}>;

export function ContentOutline({
  children,
  label,
  collectEntries,
  disabled = false,
  placement = "inline",
}: ContentOutlineProps) {
  const { rootRef, entries, activeKey, visible } = useContentOutline(
    collectEntries,
    disabled,
  );
  const styles = contentOutline({ placement });

  return (
    <div
      className={styles.root}
      data-visible={visible}
      data-placement={placement}
    >
      <div className={styles.body} ref={rootRef}>
        {children}
      </div>

      {visible && (
        <ContentOutlineRail
          placement={placement}
          label={label}
          entries={entries}
          activeKey={activeKey}
        />
      )}
    </div>
  );
}
