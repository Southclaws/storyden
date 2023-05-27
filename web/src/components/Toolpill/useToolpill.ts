import { useOutsideClick } from "@chakra-ui/react";
import { RefObject, useRef } from "react";

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
