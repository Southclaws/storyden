import { useRouter } from "next/router";
import { ReactElement } from "react";
import { Fullpage } from "src/layouts/Fullpage";
import { AuthScreen } from "src/screens/auth/AuthScreen";
import { z } from "zod";

const QuerySchema = z.object({ method: z.string().optional() });

function Page() {
  const { query } = useRouter();
  const { method } = QuerySchema.parse(query);

  return <AuthScreen method={method} />;
}

Page.getLayout = function getLayout(page: ReactElement) {
  return <Fullpage>{page}</Fullpage>;
};

export default Page;
