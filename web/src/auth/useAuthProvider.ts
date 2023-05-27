import { useRouter } from "next/router";
import { useEffect } from "react";
import { useAccountGet } from "src/api/openapi/accounts";
import { Account } from "src/api/openapi/schemas";

const PRIVATE_PAGES = ["/settings", "/new"];

function privatePage(pathName: string): boolean {
  return PRIVATE_PAGES.includes(pathName);
}

type UseAuthProvider = {
  firstTime: boolean;
  account: Account | undefined;
};

export function useAuthProvider(): UseAuthProvider {
  const { push, pathname } = useRouter();
  const { isLoading, data, error } = useAccountGet();

  const loggedIn = Boolean(data) && !error;
  const firstTime = data === undefined && error === undefined;
  const isPrivate = pathname && privatePage(pathname);

  useEffect(() => {
    if (isLoading) return;

    if (!loggedIn && isPrivate) {
      console.log("redirecting to /auth");
      push("/auth");
    }
    if (loggedIn && pathname === "/auth") {
      console.log("redirecting to /");
      push("/");
    }
  }, [isLoading, loggedIn, isPrivate, pathname, push]);

  return { firstTime, account: data };
}
