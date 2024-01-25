"use client";

import { z } from "zod";

import { ComposeScreen } from "src/screens/compose/ComposeScreen";

const QuerySchema = z.object({
  id: z.string().optional(),
});

type Props = {
  searchParams: z.infer<typeof QuerySchema>;
};

export default function Page(props: Props) {
  const params = QuerySchema.parse(props.searchParams);
  return <ComposeScreen editing={params.id} />;
}
