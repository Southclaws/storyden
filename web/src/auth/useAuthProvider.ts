import { useRouter } from "next/router";
import { useEffect } from "react";
import { useAccountGet } from "src/api/openapi/accounts";
import { Account } from "src/api/openapi/schemas";

const PRIVATE_PAGES = ["/settings"];

function privatePage(pathName: string): boolean {
  return PRIVATE_PAGES.includes(pathName);
}

type UseAuthProvider = {
  firstTime: boolean;
  account: Account | undefined;
};

export function useAuthProvider(): UseAuthProvider {
  const { push, pathname } = useRouter();
  const account = useAccountGet();

  const loggedIn = Boolean(account.data) && !account.error;
  const firstTime = account.data === undefined && account.error === undefined;
  const isPrivate = pathname && privatePage(pathname);

  useEffect(() => {
    if (!loggedIn && isPrivate) {
      push("/auth");
    }
    if (loggedIn && pathname === "/auth") {
      push("/");
    }
  }, [loggedIn, isPrivate, pathname, push]);

  return { firstTime, account: account.data };
}
