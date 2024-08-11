"use client";

import {
  PropsWithChildren,
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { createPortal } from "react-dom";

type NavigationContextType = {
  ref: any;
  setRef: (ref: any) => void;
};

export const NavigationContext = createContext<NavigationContextType>({
  ref: null,
  setRef: () => {},
});

export function NavigationProvider({ children }: PropsWithChildren) {
  const [ref, setRef] = useState();
  return (
    <NavigationContext.Provider value={{ ref, setRef }}>
      {children}
    </NavigationContext.Provider>
  );
}

export function useRightPortalRef() {
  const { setRef, ref } = useContext(NavigationContext);
  const targetRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (targetRef.current) {
      setRef(targetRef.current);
    }
  }, [targetRef, setRef]);

  return { targetRef, ref };
}

export function RightNavPortal({ children }: PropsWithChildren) {
  const { ref } = useRightPortalRef();

  if (!ref) return null;

  return <>{createPortal(children, ref)}</>;
}
