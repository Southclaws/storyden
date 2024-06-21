import { Metadata } from "next";
import { PropsWithChildren } from "react";

import { Default } from "src/layouts/Default";
import { getInfo } from "src/utils/info";

export default async function Layout({ children }: PropsWithChildren) {
  return <Default>{children}</Default>;
}

export async function generateMetadata(): Promise<Metadata> {
  const info = await getInfo();

  return {
    title: `Draft a new post on ${info.title}`,
    description: `Compose a new masterpice and share it with the community on ${info.title}`,
  };
}
