import { useOutsideClick } from "@chakra-ui/react";
import { useRouter } from "next/router";
import { RefObject, useEffect, useRef, useState } from "react";
import { useAuthProvider } from "src/auth/useAuthProvider";

export function useNavpill() {
  const { asPath } = useRouter();
  const overlayRef = useRef<HTMLDivElement>() as RefObject<HTMLDivElement>;
  const [isExpanded, setExpanded] = useState(false);
  const { account } = useAuthProvider();

  // Close the menu for either navigation events or outside clicks/taps:

  useEffect(() => setExpanded(false), [asPath]);

  useOutsideClick({
    ref: overlayRef,
    handler: () => setExpanded(false),
  });

  const onExpand = () => {
    setExpanded(!isExpanded);
  };

  return {
    overlayRef,
    isExpanded,
    onExpand,
    account,
  };
}
