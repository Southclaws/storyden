"use client";

import Link from "next/link";
import { ReactNode, useCallback } from "react";

import { button } from "@/styled-system/recipes";
import { scrollToTop } from "@/utils/scroll";

interface ScrollToTopProps {
  children?: ReactNode;
}

export function ScrollToTop({ children = "scroll to top" }: ScrollToTopProps) {
  const handleClick = useCallback((e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault();
    scrollToTop();
  }, []);

  return (
    <Link
      href="#"
      className={button({
        variant: "subtle",
        size: "xs",
      })}
      onClick={handleClick}
    >
      {children}
    </Link>
  );
}
