import { PropsWithChildren } from "react";

import { Default } from "@/layouts/Default";

export default async function Layout({
  children,
  sidebar,
}: PropsWithChildren<{ sidebar: React.ReactNode }>) {
  return <Default sidebar={sidebar}>{children}</Default>;
}
