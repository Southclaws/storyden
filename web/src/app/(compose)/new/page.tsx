"use client";

import { useSearchParams } from "next/navigation";

import { ComposeScreen } from "src/screens/compose/ComposeScreen";

export default function Page() {
  const params = useSearchParams();

  return <ComposeScreen editing={params.get("id") ?? undefined} />;
}
