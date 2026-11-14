import type { ReactNode } from "react";

import { styled } from "@/styled-system/jsx";

type MetaGridProps = {
  children: ReactNode;
  columns?: {
    base: string;
    md: string;
  };
};

export function MetaGrid({ children, columns }: MetaGridProps) {
  return (
    <styled.dl
      display="grid"
      gridTemplateColumns={
        columns ?? {
          base: "minmax(0, 1fr)",
          md: "repeat(auto-fit, minmax(min(100%, 12rem), 1fr))",
        }
      }
      gap={{ base: "2", md: "3" }}
      color="text.subtle"
      fontSize="xs"
      w="full"
    >
      {children}
    </styled.dl>
  );
}

export function MetaItem({
  label,
  children,
}: {
  label: string;
  children: ReactNode;
}) {
  return (
    <styled.div minW="0">
      <styled.dt
        fontSize="xs"
        fontWeight="bold"
        color="text.subtle"
        textTransform="uppercase"
      >
        {label}
      </styled.dt>
      <styled.dd wordBreak="break-word">{children}</styled.dd>
    </styled.div>
  );
}
