import type { Metadata } from "next";

import { HomeScreen } from "src/screens/home/HomeScreen";
import { getInfo } from "src/utils/info";

export default async function Page() {
  return <HomeScreen />;
}

export async function generateMetadata(): Promise<Metadata> {
  const info = await getInfo();

  return {
    title: info.title,
    description: info.description,
  };
}
