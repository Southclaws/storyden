"use client";

import Link from "next/link";
import { ReactNode, useCallback } from "react";

import { css } from "@/styled-system/css";
import { scrollToTop } from "@/utils/scroll";

interface ScrollToTopProps {
  children?: ReactNode;
}

export function ScrollToTop({ 
  children = "scroll to top"
}: ScrollToTopProps) {
  const handleClick = useCallback((e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault();
    scrollToTop();
  }, []);

  return (
    <Link
      href="#"
      style={{
        border: "1px solid var(--colors-border-muted)",
        transition: "all 200ms ease-out"
      }}
      className={css({
        display: "inline-flex",
        alignItems: "center",
        justifyContent: "center",
        px: "2",
        py: "1",
        fontSize: "xs",
        fontWeight: "medium",
        borderRadius: "sm",
        textDecoration: "none",
        cursor: "pointer",
        color: "fg.muted",
        backgroundColor: "bg.subtle",
        colorPalette: "accent",
        _hover: {
          backgroundColor: "accent.subtle",
          borderColor: "accent.muted",
          color: "accent.default",
          transform: "translateY(-1px)",
          boxShadow: "sm"
        },
        _active: {
          transform: "translateY(0)",
          boxShadow: "xs"
        }
      })}
      onClick={handleClick}
    >
      {children}
    </Link>
  );
}
