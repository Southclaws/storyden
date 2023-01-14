import { GetServerSidePropsContext, GetServerSidePropsResult } from "next";
import Link from "next/link";
import { useRouter } from "next/router";
import nookies from "nookies";
import { useEffect } from "react";

export default function Page() {
  const router = useRouter();
  useEffect(() => {
    document.cookie =
      "storyden-session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";

    router.push("/");
  }, [router]);

  return <Link href="/login">Logged out. Returning to login.</Link>;
}

export async function getServerSideProps(
  ctx: GetServerSidePropsContext
): Promise<GetServerSidePropsResult<{}>> {
  ctx.res.setHeader("Clear-Site-Data", "cookies");
  ctx.res.setHeader("Set-Cookie", "storyden-session=x; Max-Age=0");
  nookies.destroy(ctx, "storyden-session");
  return {
    props: {},
  };
}
