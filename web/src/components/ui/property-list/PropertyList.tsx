import { Fragment, ReactNode } from "react";

import { styled } from "@/styled-system/jsx";
import { HTMLStyledProps } from "@/styled-system/types";

export type PropertyListItem = {
  id: string;
  label: ReactNode;
  value: ReactNode;
};

export type PropertyListProps = Omit<HTMLStyledProps<"dl">, "children"> & {
  items: PropertyListItem[];
};

export function PropertyList({ items, ...props }: PropertyListProps) {
  if (items.length === 0) {
    return null;
  }

  return (
    <styled.dl
      w="full"
      minW="0"
      fontSize="xs"
      display="grid"
      gridTemplateColumns="fit-content(40%) minmax(0, 1fr)"
      alignItems="center"
      columnGap="2"
      {...props}
    >
      {items.map(({ id, label, value }) => (
        <Fragment key={id}>
          <styled.dt minW="0" color="text.muted" overflowWrap="anywhere">
            {label}
          </styled.dt>
          <styled.dd
            minW="0"
            textAlign="right"
            overflowWrap="anywhere"
            lineClamp={1}
          >
            {value}
          </styled.dd>
        </Fragment>
      ))}
    </styled.dl>
  );
}
