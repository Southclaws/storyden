"use client";

import { RefObject, useRef } from "react";

import { useOutsideClick } from "src/theme/components";

export type Props = {
  onClickOutside?: () => void;
};

export function useToolpill({ onClickOutside }: Props) {
  const ref = useRef<HTMLDivElement>() as RefObject<HTMLDivElement>;

  useOutsideClick({
    ref: ref,
    handler: () => {
      onClickOutside?.();
    },
  });

  return { ref };
}
