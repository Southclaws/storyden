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
          base: "1fr",
          md: "minmax(12rem, 1fr) minmax(12rem, 1fr) minmax(8rem, .7fr)",
        }
      }
      gap={{ base: "2", md: "3" }}
      color="fg.muted"
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
        color="fg.muted"
        textTransform="uppercase"
      >
        {label}
      </styled.dt>
      <styled.dd wordBreak="break-word">{children}</styled.dd>
    </styled.div>
  );
}
