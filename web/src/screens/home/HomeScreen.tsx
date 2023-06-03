"use client";

import { Unready } from "src/components/Unready";
import { ThreadList } from "./components/ThreadList";
import { useHomeScreen } from "./useHomeScreen";

export function HomeScreen() {
  const { data, error } = useHomeScreen();

  if (!data) return <Unready {...error} />;

  return <ThreadList threads={data.threads} />;
}
