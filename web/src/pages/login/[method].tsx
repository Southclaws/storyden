import { NextPage } from "next";
import { useRouter } from "next/router";
import { AuthScreen } from "../../screens/auth/AuthScreen";

const Page: NextPage = () => {
  const { query } = useRouter();
  const method = query["method"] as string;
  return <AuthScreen method={method} />;
};

export default Page;
