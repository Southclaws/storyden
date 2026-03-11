import { PropsWithChildren } from "react";

export default async function Layout({ children }: PropsWithChildren) {
  return <div className="sd-screen sd-screen--discussion">{children}</div>;
}
