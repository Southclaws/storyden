"use client";

import { useParams } from "next/navigation";

import { TrailEditScreen } from "@/screens/trails/TrailEditScreen";

export default function Page() {
  const { id } = useParams<{ id: string }>();
  return <TrailEditScreen trailId={id} />;
}
