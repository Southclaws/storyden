"use client";

import { useParams } from "next/navigation";
import { z } from "zod";

import { Unready } from "src/components/site/Unready";
import { CollectionScreen } from "src/screens/collection/CollectionScreen";

const ParamSchema = z.object({
  handle: z.string().optional(),
  collection: z.string().optional(),
});

export default function Page() {
  const params = useParams();
  const { handle, collection } = ParamSchema.parse(params);

  if (!handle || !collection) return <Unready />;

  return <CollectionScreen handle={handle} collection={collection} />;
}
