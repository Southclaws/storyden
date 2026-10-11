"use client";

import { Portal } from "@ark-ui/react";
import { type CSSProperties, Fragment } from "react";

import { Button } from "@/components/ui/button";
import * as Popover from "@/components/ui/popover";
import { RelativeTime } from "@/components/ui/relative-time";
import { contentOutline } from "@/styled-system/recipes";

import type { ContentOutlineEntry } from "./outline";
import { useContentOutlineRail } from "./useContentOutlineRail";

type Props = {
  label: string;
  entries: ContentOutlineEntry[];
  activeKey?: string;
};

export function ContentOutlineRail({ label, entries, activeKey }: Props) {
  const styles = contentOutline();
  const { open, triggerRef, panelRef, ticksRef, menuScrollRef, handlers } =
    useContentOutlineRail(activeKey);

  const items = entries.map((entry, index) => {
    const previousEntry = entries[index - 1];
    const isReply = entry.kind === "reply";
    const followsContent = previousEntry?.kind !== "reply";

    return {
      entry,
      active: entry.key === activeKey,
      startsReplies: isReply && followsContent,
    };
  });

  return (
    <aside className={styles.rail} data-content-outline="">
      <Popover.Root
        open={open}
        autoFocus={false}
        restoreFocus={false}
        modal={false}
        positioning={{ placement: "left", gutter: 8 }}
        onOpenChange={handlers.handleOpenChange}
        onFocusOutside={handlers.handleFocusOutside}
        onEscapeKeyDown={handlers.handleEscapeKeyDown}
      >
        <Popover.Trigger asChild>
          <Button
            className={styles.trigger}
            ref={triggerRef}
            variant="plain"
            aria-label={label}
            onPointerEnter={handlers.handleTriggerPointerEnter}
            onPointerLeave={handlers.handleCloseSoon}
            onFocus={handlers.handleTriggerFocus}
            onBlur={handlers.handleTriggerBlur}
            onClick={handlers.handleTriggerClick}
            onKeyDown={handlers.handleTriggerKeyDown}
          >
            <span className={styles.ticks} ref={ticksRef} aria-hidden="true">
              {items.map(({ entry, active, startsReplies }) => (
                <span
                  className={styles.tick}
                  key={entry.key}
                  data-kind={entry.kind}
                  data-active={active}
                  data-reply-start={startsReplies}
                  style={{ "--outline-depth": entry.depth } as CSSProperties}
                />
              ))}
            </span>
          </Button>
        </Popover.Trigger>

        <Portal>
          <Popover.Positioner>
            <Popover.Content
              className={styles.content}
              ref={panelRef}
              onPointerEnter={handlers.handleKeepOpen}
              onPointerLeave={handlers.handleCloseSoon}
              onFocus={handlers.handlePanelFocus}
              onBlur={handlers.handleCloseSoon}
            >
              <div className={styles.viewport} ref={menuScrollRef}>
                <Popover.Title>{label}</Popover.Title>
                <nav className={styles.navigation} aria-label={label}>
                  {items.map(({ entry, active, startsReplies }) => (
                    <Fragment key={entry.key}>
                      {startsReplies && (
                        <p className={styles.replyLabel}>Replies</p>
                      )}
                      <ContentOutlineItem
                        entry={entry}
                        active={active}
                        onNavigate={handlers.handleNavigate}
                      />
                    </Fragment>
                  ))}
                </nav>
              </div>
            </Popover.Content>
          </Popover.Positioner>
        </Portal>
      </Popover.Root>
    </aside>
  );
}

type ItemProps = {
  entry: ContentOutlineEntry;
  active: boolean;
  onNavigate: (entry: ContentOutlineEntry) => void;
};

function ContentOutlineItem({ entry, active, onNavigate }: ItemProps) {
  const styles = contentOutline();

  return (
    <Button
      className={styles.item}
      variant="ghost"
      aria-current={active ? "location" : undefined}
      data-kind={entry.kind}
      style={{ "--outline-depth": entry.depth } as CSSProperties}
      title={entry.label}
      onClick={() => onNavigate(entry)}
    >
      <span className={styles.itemContent}>
        <span className={styles.byline}>
          <span className={styles.label} data-part="label">
            {entry.label}
          </span>
          {entry.createdAt && (
            <RelativeTime
              className={styles.timestamp}
              value={entry.createdAt}
              title={entry.createdAt}
            />
          )}
        </span>
        {entry.description && (
          <span className={styles.snippet}>{entry.description}</span>
        )}
      </span>
    </Button>
  );
}
