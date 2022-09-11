import {
  GetServerSidePropsContext,
  GetServerSidePropsResult,
  NextPage,
} from "next";
import { useRouter } from "next/router";
import { AuthScreen } from "../../screens/auth/AuthScreen";
import * as z from "zod";

const ParamsSchema = z.object({
  method: z.string(),
});

const Page: NextPage = () => {
  const { query } = useRouter();
  const method = query.method as string;
  return <AuthScreen method={method} />;
};

export default Page;
