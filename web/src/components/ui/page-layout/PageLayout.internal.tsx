import type { ComponentProps } from "react";
import { forwardRef } from "react";

import { cx } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

export type PageLayoutWidth = "content" | "wide" | "full";

const StyledPageLayout = styled("div", {
  base: {
    width: "full",
    minWidth: "0",
  },
});

export interface PageLayoutProps extends Omit<
  ComponentProps<typeof StyledPageLayout>,
  "width"
> {
  width?: PageLayoutWidth;
}

export const PageLayout = forwardRef<HTMLDivElement, PageLayoutProps>(
  function PageLayout({ width = "content", className, ...props }, ref) {
    return (
      <StyledPageLayout
        {...props}
        ref={ref}
        className={cx("page-layout", className)}
        data-page-width={width}
      />
    );
  },
);
