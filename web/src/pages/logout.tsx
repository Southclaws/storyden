import { GetServerSidePropsContext, GetServerSidePropsResult } from "next";
import Link from "next/link";
import nookies from "nookies";

export default function Page() {
  document.cookie =
    "storyden-session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
  return <Link href="/login">Logged out. Returning to login.</Link>;
}

export async function getServerSideProps(
  ctx: GetServerSidePropsContext
): Promise<GetServerSidePropsResult<{}>> {
  ctx.res.setHeader("Clear-Site-Data", "cookies");
  nookies.destroy(ctx, "storyden-session");

  return {
    props: {},
    redirect: {
      destination: "/",
      permanent: false,
    },
  };
}
