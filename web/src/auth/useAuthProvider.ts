import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";

import { useAccountGet } from "src/api/openapi/accounts";

const PRIVATE_PAGES = ["/settings", "/new", "/admin"];

function privatePage(pathName: string): boolean {
  return PRIVATE_PAGES.includes(pathName);
}

export function useAuthProvider() {
  const { isLoading, data, error } = useAccountGet();
  const { push } = useRouter();
  const pathname = usePathname();

  const loggedIn = Boolean(data) && !error;
  const isPrivate = pathname && privatePage(pathname);

  useEffect(() => {
    if (isLoading) return;

    if (!loggedIn && isPrivate) {
      push("/auth");
    }
    if (loggedIn && pathname === "/auth") {
      push("/");
    }
  }, [isLoading, loggedIn, isPrivate, pathname, push]);

  return;
}
