"use client";

import { useParams } from "next/navigation";

import { TrailDetailScreen } from "@/screens/trails/TrailDetailScreen";

export default function Page() {
  const { id } = useParams<{ id: string }>();
  return <TrailDetailScreen trailId={id} />;
}
