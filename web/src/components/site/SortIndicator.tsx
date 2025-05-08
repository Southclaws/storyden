import { useState } from "react";

import { ChevronDownIcon, ChevronUpIcon } from "@/components/ui/icons/Chevron";
import { Box, styled } from "@/styled-system/jsx";

type Direction = "asc" | "desc" | "none";

type SortState = {
  property: string;
  order: Direction;
};

type Props = {
  order: Direction;
};

export function useSortIndicator() {
  const [sort, setSort] = useState<SortState | null>(null);

  function handleSort(property: string) {
    // cycle through sort states: none -> asc -> desc -> none
    if (sort?.property === property) {
      if (sort.order === "asc") {
        setSort({ property, order: "desc" });
      } else if (sort.order === "desc") {
        setSort(null);
      }
    } else {
      setSort({ property, order: "asc" });
    }
  }

  return {
    sort,
    handleSort,
  };
}

export function SortIndicator({ order }: Props) {
  const label =
    order === "none"
      ? "No sort"
      : order === "asc"
        ? "Sort ascending"
        : "Sort descending";

  return (
    <styled.span width="4" height="4" aria-label={label} title={label}>
      {order === "none" ? (
        <>&nbsp;</>
      ) : order === "asc" ? (
        <ChevronUpIcon width="4" height="4" />
      ) : (
        <ChevronDownIcon width="4" height="4" />
      )}
    </styled.span>
  );
}
