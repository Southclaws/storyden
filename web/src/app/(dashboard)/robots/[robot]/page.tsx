"use client";

import { RobotConfigurationScreen } from "@/screens/robots/RobotConfigurationScreen";

type Props = {
  params: Promise<{
    robot: string;
  }>;
};

export default async function Page(props: Props) {
  const { robot } = await props.params;

  return <RobotConfigurationScreen robotId={robot} />;
}
