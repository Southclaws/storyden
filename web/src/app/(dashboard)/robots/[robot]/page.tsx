"use client";

import { useParams } from "next/navigation";

import { RobotConfigurationScreen } from "@/screens/robots/RobotConfigurationScreen";

export default function Page() {
  const { robot } = useParams<{ robot: string }>();

  return <RobotConfigurationScreen robotId={robot} />;
}
