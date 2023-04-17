import { useOutsideClick } from "@chakra-ui/react";
import { RefObject, useRef, useState } from "react";
import { useAuthProvider } from "src/auth/useAuthProvider";

export function useNavpill() {
  const overlayRef = useRef<HTMLDivElement>() as RefObject<HTMLDivElement>;
  const [isExpanded, setExpanded] = useState(false);
  const { account } = useAuthProvider();

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
    isAuthenticated: !!account,
  };
}
