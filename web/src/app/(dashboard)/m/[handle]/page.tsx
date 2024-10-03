"use client";

import { useParams } from "next/navigation";

import { Unready } from "src/components/site/Unready";
import { ProfileScreen } from "src/screens/profile/ProfileScreen";

import { ParamSchema } from "./utils";

export default function Page() {
  const params = useParams();
  const { handle } = ParamSchema.parse(params);

  if (!handle) return <Unready />;

  return <ProfileScreen handle={handle} />;
}
