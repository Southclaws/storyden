"use client";

import { useState } from "react";

import { useRobotsList } from "@/api/openapi-client/robots";
import { Robot } from "@/api/openapi-schema";
import { RobotCard } from "@/components/robots/RobotCard";
import { RobotConfigurationModal } from "@/components/robots/RobotConfiguration/RobotConfigurationModal";
import { Unready } from "@/components/site/Unready";
import { LStack } from "@/styled-system/jsx";
import { useDisclosure } from "@/utils/useDisclosure";

import { useSelectedRobot } from "./useSelectedRobot";

export default function RobotList() {
  const { data, error } = useRobotsList();
  const disclosure = useDisclosure();
  const [robotID, setSelectedRobot] = useSelectedRobot();

  if (!data) {
    return <Unready error={error} />;
  }

  function handleRobotClick(robot: string) {
    setSelectedRobot(robot);
    disclosure.onOpen();
  }

  function handleClose() {
    disclosure.onClose();
    setSelectedRobot(null);
  }

  return (
    <LStack>
      {data.robots.map((robot) => (
        <RobotCard
          key={robot.id}
          robot={robot}
          onClick={() => handleRobotClick(robot.id)}
        />
      ))}
    </LStack>
  );
}
