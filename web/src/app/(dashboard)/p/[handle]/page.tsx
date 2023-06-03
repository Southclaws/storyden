"use client";

import { useParams } from "next/navigation";
import { Unready } from "src/components/Unready";
import { ProfileScreen } from "src/screens/profile/ProfileScreen";
import { z } from "zod";

const ParamSchema = z.object({ handle: z.string().optional() });

export default function Page() {
  const params = useParams();
  const { handle } = ParamSchema.parse(params);

  if (!handle) return <Unready />;

  return <ProfileScreen handle={handle} />;
}
