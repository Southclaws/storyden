import { useState } from "react";
import { useAuthProvider } from "src/auth/useAuthProvider";

export function useNavigation() {
  const { account } = useAuthProvider();
  const [isExpanded, setExpanded] = useState(false);

  const onExpand = () => {
    setExpanded(!isExpanded);
  };

  return {
    isAuthenticated: !!account,
    isExpanded,
    onExpand,
  };
}
