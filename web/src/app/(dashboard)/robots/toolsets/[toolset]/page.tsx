"use client";

import { useParams } from "next/navigation";

import { RobotToolsetConfigurationScreen } from "@/screens/robots/RobotToolsetConfigurationScreen";

export default function Page() {
  const { toolset } = useParams<{ toolset: string }>();
  return <RobotToolsetConfigurationScreen toolsetId={toolset} />;
}
