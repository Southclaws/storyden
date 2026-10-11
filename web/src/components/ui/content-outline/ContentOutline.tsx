"use client";

import type { PropsWithChildren } from "react";

import { contentOutline } from "@/styled-system/recipes";

import { ContentOutlineRail } from "./ContentOutline.internal";
import type { CollectOutlineEntries } from "./outline";
import { useContentOutline } from "./useContentOutline";

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
  const { rootRef, entries, activeKey } = useContentOutline(
    collectEntries,
    disabled,
  );
  const styles = contentOutline();
  const visible = entries.length > 1;

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
