import { ReactElement } from "react";
import { Fullpage } from "src/layouts/Fullpage";
import { AuthScreen } from "src/screens/auth/AuthScreen";

function Page() {
  return <AuthScreen method={null} />;
}

Page.getLayout = function getLayout(page: ReactElement) {
  return <Fullpage>{page}</Fullpage>;
};

export default Page;
