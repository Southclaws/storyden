import { PropsWithChildren } from "react";

import { Default } from "@/layouts/Default";

export default async function Layout({
  children,
  contextpane,
}: PropsWithChildren<{ contextpane: React.ReactNode }>) {
  return <Default contextpane={contextpane}>{children}</Default>;
}
