import { GetServerSidePropsContext, GetServerSidePropsResult } from "next";
import Link from "next/link";
import { useRouter } from "next/router";
import { destroyCookie } from "nookies";
import { useEffect } from "react";
import { getAccountGetKey } from "src/api/openapi/accounts";
import { useAuthProviderLogout } from "src/api/openapi/auth";
import { useSWRConfig } from "swr";

export default function Page() {
  const router = useRouter();
  useAuthProviderLogout();
  const { mutate } = useSWRConfig();

  useEffect(() => {
    mutate(getAccountGetKey(), null);
    router.push("/");
  }, [router, mutate]);

  return <Link href="/login">Logged out. Returning to login.</Link>;
}

export async function getServerSideProps(
  ctx: GetServerSidePropsContext
): Promise<GetServerSidePropsResult<object>> {
  ctx.res.setHeader("Clear-Site-Data", `"cookies"`);
  destroyCookie(ctx, "storyden-session");

  return {
    props: {},
  };
}
