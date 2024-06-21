import { PropsWithChildren } from "react";

import { Default } from "src/layouts/Default";

export default async function Layout({ children }: PropsWithChildren) {
  return <Default>{children}</Default>;
}
