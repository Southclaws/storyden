import chroma from "chroma-js";
import Link from "next/link";

import { TagReference } from "@/api/openapi-schema";
import { css, cx } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";
import { badge } from "@/styled-system/recipes";

import { BadgeProps } from "../ui/badge";

export type Props = BadgeProps & {
  tag: TagReference;
  showItemCount?: boolean;
};

// tags are always lowercase, which means most ascenders and descenders are
// slightly mis-aligned to the optical center of the badge.
const OPTICAL_ALIGNMENT_ADJUSTMENT = 0.5;

const badgeStyles = css({
  bgColor: "colorPalette",
  borderColor: "colorPalette.muted",
  color: "colorPalette.text",
});

export function TagBadge({ tag, showItemCount, ...props }: Props) {
  const cssVars = badgeColourCSS(tag.colour);

  const styles = {
    ...cssVars,
    "--optical-adjustment-top": `${-OPTICAL_ALIGNMENT_ADJUSTMENT}px`,
    "--optical-adjustment-bot": `${OPTICAL_ALIGNMENT_ADJUSTMENT}px`,
    "--optical-adjustment-count-right": "0.4rem",
  };

  const titleLabel = `${tag.item_count} items tagged with ${tag.name}`;

  return (
    <Link
      className={cx(
        badge({
          size: "sm",
          ...props,
        }),
        badgeStyles,
      )}
      style={styles}
      title={titleLabel}
      href={`/tags/${tag.name}`}
    >
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
    </Link>
  );
}

function badgeColourCSS(c: string) {
  const { bg, bo, fg } = badgeColours(c);

  return {
    "--colors-color-palette-text": fg,
    "--colors-color-palette-muted": bo,
    "--colors-color-palette": bg,
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
