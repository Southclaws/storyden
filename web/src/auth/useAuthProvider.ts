import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";

import { useAccountGet } from "@/api/openapi-client/accounts";

import { usePublicRegistration } from "@/lib/settings/registration";

const PRIVATE_PAGES = ["/settings", "/new", "/admin"];

function privatePage(pathName: string): boolean {
  return PRIVATE_PAGES.includes(pathName);
}

export function useAuthProvider() {
  const { isLoading, data, error } = useAccountGet();
  const { push } = useRouter();
  const pathname = usePathname();
  const canRegister = usePublicRegistration();

  const loggedIn = Boolean(data) && !error;
  const isPrivate = pathname && privatePage(pathname);

  useEffect(() => {
    if (isLoading) return;

    if (!loggedIn && isPrivate) {
      push(canRegister ? "/register" : "/login");
    }
    if (loggedIn && (pathname === "/register" || pathname === "/login")) {
      push("/");
    }
  }, [canRegister, isLoading, loggedIn, isPrivate, pathname, push]);

  return;
}
