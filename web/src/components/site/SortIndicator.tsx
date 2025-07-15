import { useState } from "react";

import {
  ChevronDownIcon,
  ChevronUpDownIcon,
  ChevronUpIcon,
} from "@/components/ui/icons/Chevron";
import { styled } from "@/styled-system/jsx";

type Direction = "asc" | "desc" | "none";

export type SortState = {
  property: string;
  order: Direction;
};

export type SortIndicatorProps = {
  order: Direction;
};

export type UseSortIndicator = ReturnType<typeof useSortIndicator>;

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

export function SortIndicator({ order }: SortIndicatorProps) {
  const label =
    order === "none"
      ? "No sort"
      : order === "asc"
        ? "Sort ascending"
        : "Sort descending";

  return (
    <styled.span width="4" height="4" aria-label={label} title={label}>
      {order === "none" ? (
        <ChevronUpDownIcon width="4" height="4" color="fg.muted" />
      ) : order === "asc" ? (
        <ChevronUpIcon width="4" height="4" />
      ) : (
        <ChevronDownIcon width="4" height="4" />
      )}
    </styled.span>
  );
}
