"use client";

import { useClickAway } from "@uidotdev/usehooks";

export type Props = {
  onClickOutside?: () => void;
};

export function useToolpill({ onClickOutside }: Props) {
  const ref = useClickAway<HTMLDivElement>(() => {
    onClickOutside?.();
  });

  return { ref };
}
