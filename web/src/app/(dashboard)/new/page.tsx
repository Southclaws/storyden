"use client";

import { use } from "react";
import { z } from "zod";

import { ComposeScreen } from "src/screens/compose/ComposeScreen";

const QuerySchema = z.object({
  id: z.string().optional(),
});

type Props = {
  searchParams: Promise<z.infer<typeof QuerySchema>>;
};

export default function Page(props: Props) {
  const searchParams = use(props.searchParams);
  const params = QuerySchema.parse(searchParams);
  return <ComposeScreen editing={params.id} />;
}
