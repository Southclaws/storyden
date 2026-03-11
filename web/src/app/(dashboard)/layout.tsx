import { PropsWithChildren } from "react";

import { Default } from "src/layouts/Default";

export default async function Layout({
  children,
  contextpane,
}: PropsWithChildren<{ contextpane: React.ReactNode }>) {
  return (
    <Default contextpane={contextpane}>
      <div className="sd-screen sd-screen--dashboard-route">{children}</div>
    </Default>
  );
}
