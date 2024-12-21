import { useCallback } from "react";

export function useMergedRefs<T>(...refs: (React.Ref<T> | undefined | null)[]) {
  return useCallback(
    (node: T | null) => {
      refs.forEach((ref) => {
        if (typeof ref === "function") {
          ref(node);
        } else if (ref && typeof ref === "object") {
          (ref as React.MutableRefObject<T | null>).current = node;
        }
      });
    },
    [refs],
  );
}
