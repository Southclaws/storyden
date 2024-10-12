import { useEffect, useState } from "react";

export function useClientState<T>(fn: () => T | null) {
  const [state, setState] = useState<T | null>(null);

  useEffect(() => {
    const value = fn();

    setState(() => value);
  }, []);

  return state;
}

export function useShare() {
  return useClientState(() => {
    if (navigator.share) {
      return navigator.share;
    }

    return null;
  });
}
