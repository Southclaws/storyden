import { formatDistanceToNow, formatDistanceToNowStrict } from "date-fns";

import { styled } from "@/styled-system/jsx";
import type { HTMLStyledProps } from "@/styled-system/types";

type RelativeTimeProps = Omit<HTMLStyledProps<"time">, "children"> & {
  precision?: "approximate" | "strict";
  value: string | number | Date;
};

export function RelativeTime({
  precision = "approximate",
  value,
  ...props
}: RelativeTimeProps) {
  const date = value instanceof Date ? value : new Date(value);
  const label =
    precision === "strict"
      ? formatDistanceToNowStrict(date, { addSuffix: true })
      : formatDistanceToNow(date, { addSuffix: true });

  return (
    <styled.time
      dateTime={date.toISOString()}
      suppressHydrationWarning
      {...props}
    >
      {label}
    </styled.time>
  );
}
