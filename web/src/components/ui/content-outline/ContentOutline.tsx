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
}>;

export function ContentOutline({
  children,
  label,
  collectEntries,
  disabled = false,
}: ContentOutlineProps) {
  const { rootRef, entries, activeKey, visible } = useContentOutline(
    collectEntries,
    disabled,
  );
  const styles = contentOutline();

  return (
    <div className={styles.root} data-visible={visible}>
      <div className={styles.body} ref={rootRef}>
        {children}
      </div>

      {visible && (
        <ContentOutlineRail
          label={label}
          entries={entries}
          activeKey={activeKey}
        />
      )}
    </div>
  );
}
