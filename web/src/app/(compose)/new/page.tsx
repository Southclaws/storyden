"use client";

import { z } from "zod";
import dynamic from "next/dynamic";
import { useSearchParams } from "next/navigation";

// NOTE: Render client side, probably unnecessary but quick fix to react errors.
const ComposeScreen = dynamic(
  () =>
    import("../../../screens/compose/ComposeScreen").then(
      (mod) => mod.ComposeScreen
    ),
  {
    ssr: false,
  }
);

const QuerySchema = z.object({
  id: z.string().optional(),
});

export default function Page() {
  const params = useSearchParams();
  const { id } = QuerySchema.parse(params);

  return <ComposeScreen editing={id} />;
}
