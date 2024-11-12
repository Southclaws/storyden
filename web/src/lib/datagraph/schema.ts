import { z } from "zod";

import { DatagraphItemKind } from "@/api/openapi-schema";

export const DatagraphKindTable: { [K in DatagraphItemKind]: string } = {
  post: "Post",
  thread: "Thread",
  reply: "Reply",
  node: "Library",
  collection: "Collection",
  profile: "Profile",
  event: "Event",
};

const DatagraphKinds = Object.keys(DatagraphKindTable) as unknown as readonly [
  DatagraphItemKind,
  ...DatagraphItemKind[],
];

export const DatagraphKindSchema: z.ZodType<DatagraphItemKind> =
  z.enum(DatagraphKinds);
export type DatagraphKind = z.infer<typeof DatagraphKindSchema>;
