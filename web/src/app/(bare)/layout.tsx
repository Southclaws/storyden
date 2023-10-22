import { PropsWithChildren } from "react";

import { Fullpage } from "src/layouts/Fullpage";

export default function Layout({ children }: PropsWithChildren) {
  return <Fullpage>{children}</Fullpage>;
}
