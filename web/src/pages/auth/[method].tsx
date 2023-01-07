import { useRouter } from "next/router";
import { ReactElement } from "react";
import { Fullpage } from "src/layouts/Fullpage";
import { AuthScreen } from "src/screens/auth/AuthScreen";

function Page() {
  const { query } = useRouter();
  const method = query["method"] as string;
  return <AuthScreen method={method} />;
}

Page.getLayout = function getLayout(page: ReactElement) {
  return <Fullpage>{page}</Fullpage>;
};

export default Page;
