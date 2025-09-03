import chroma from "chroma-js";
import Link from "next/link";

import { TagReference } from "@/api/openapi-schema";
import { css, cx } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";
import { badge } from "@/styled-system/recipes";

import { BadgeProps } from "../ui/badge";

type TagBadgeProps =
  | {
      type?: "link";
      onClick?: never;
      highlighted?: never;
    }
  | {
      type: "button";
      onClick: () => void;
      highlighted?: boolean;
    };

export type Props = BadgeProps &
  TagBadgeProps & {
    tag: TagReference;
    showItemCount?: boolean;
  };

// tags are always lowercase, which means most ascenders and descenders are
// slightly mis-aligned to the optical center of the badge.
const OPTICAL_ALIGNMENT_ADJUSTMENT = 0.5;

const badgeStyles = css({
  bgColor: "colorPalette.bg",
  borderColor: "colorPalette.border",
  color: "colorPalette.fg",
});

export function TagBadge({
  type,
  onClick,
  highlighted,
  tag,
  showItemCount,
  ...props
}: Props) {
  const cssVars = badgeColourCSS(tag.colour);

  const styles = {
    ...cssVars,
    "--optical-adjustment-top": `${-OPTICAL_ALIGNMENT_ADJUSTMENT}px`,
    "--optical-adjustment-bot": `${OPTICAL_ALIGNMENT_ADJUSTMENT}px`,
    "--optical-adjustment-count-right": "0.4rem",
  };

  const titleLabel = `${tag.item_count} items tagged with ${tag.name}`;

  const render = (
    <>
      {showItemCount && (
        <styled.span
          borderRightStyle="solid"
          borderRightWidth="thin"
          borderRightColor="colorPalette.muted"
          paddingRight="var(--optical-adjustment-count-right)"
        >
          {tag.item_count}
        </styled.span>
      )}

      <styled.span
        mb="var(--optical-adjustment-bot)"
        mt="var(--optical-adjustment-top)"
      >
        {tag.name}
      </styled.span>
    </>
  );

  const tagBadgeStyles = cx(
    badge({
      size: "sm",
      ...props,
    }),
    badgeStyles,
  );

  if (type === "button") {
    const shouldShowHighlightStyles = highlighted !== undefined;

    const highlightStyles = shouldShowHighlightStyles
      ? highlighted
        ? css({
            opacity: "full",
          })
        : css({
            opacity: "5",
          })
      : undefined;

    return (
      <styled.button
        type="button"
        className={cx(tagBadgeStyles, highlightStyles)}
        style={styles}
        title={`include ${tag.name} in filter`}
        onClick={onClick}
        cursor="pointer"
        aria-pressed={!!highlighted}
      >
        {render}
      </styled.button>
    );
  }

  return (
    <Link
      className={tagBadgeStyles}
      style={styles}
      title={titleLabel}
      href={`/tags/${tag.name}`}
    >
      {render}
    </Link>
  );
}

function badgeColourCSS(c: string) {
  const { bg, bo, fg } = badgeColours(c);

  return {
    "--colors-color-palette-fg": fg,
    "--colors-color-palette-border": bo,
    "--colors-color-palette-bg": bg,
  } as React.CSSProperties;
}

function badgeColours(c: string) {
  const colour = chroma(c);

  const hue = colour.lch()[2];

  const bg = chroma(0.95, 0.1, hue, "oklch").css();
  const bo = chroma(0.85, 0.2, hue, "oklch").css();
  const fg = chroma(0.55, 0.2, hue, "oklch").css();

  return { bg, bo, fg };
}
