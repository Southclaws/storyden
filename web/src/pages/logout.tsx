import { GetServerSidePropsContext, GetServerSidePropsResult } from "next";
import Link from "next/link";

export default function Page() {
  return <Link href="/login">Logged out. Returning to login.</Link>;
}

export async function getServerSideProps(
  ctx: GetServerSidePropsContext
): Promise<GetServerSidePropsResult<{}>> {
  ctx.res.setHeader("Clear-Site-Data", `"cookies"`);

  return {
    props: {},
    redirect: {
      destination: "/",
      permanent: false,
    },
  };
}
