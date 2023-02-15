import { GetServerSidePropsContext, GetServerSidePropsResult } from "next";
import Link from "next/link";
import { useRouter } from "next/router";
import { destroyCookie } from "nookies";
import { useEffect } from "react";
import { useAuthProviderLogout } from "src/api/openapi/auth";

export default function Page() {
  const router = useRouter();
  const logout = useAuthProviderLogout();

  useEffect(() => {
    console.log(logout);
    router.push("/");
  }, [router, logout]);

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
