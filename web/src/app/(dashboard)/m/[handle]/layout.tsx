"use client";

import { useParams } from "next/navigation";
import { PropsWithChildren } from "react";

import { Unready } from "src/components/site/Unready";
import { ProfileLayout } from "src/screens/profile/ProfileLayout";

import { ParamSchema } from "./utils";

export default function Layout({ children }: PropsWithChildren) {
  const params = useParams();
  const { handle } = ParamSchema.parse(params);

  if (!handle) return <Unready />;

  return <ProfileLayout handle={handle}>{children}</ProfileLayout>;
}
