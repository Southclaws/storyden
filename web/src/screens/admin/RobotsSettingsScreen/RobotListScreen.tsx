"use client";

import { useState } from "react";

import { useRobotsList } from "@/api/openapi-client/robots";
import { Robot } from "@/api/openapi-schema";
import { RobotCard } from "@/components/robots/RobotCard";
import { RobotConfigurationModal } from "@/components/robots/RobotConfiguration/RobotConfigurationModal";
import { Unready } from "@/components/site/Unready";
import { LStack } from "@/styled-system/jsx";
import { useDisclosure } from "@/utils/useDisclosure";

export default function RobotListScreen() {
  const { data, error } = useRobotsList();
  const [selectedRobot, setSelectedRobot] = useState<Robot | null>(null);
  const disclosure = useDisclosure();

  if (!data) {
    return <Unready error={error} />;
  }

  function handleRobotClick(robot: Robot) {
    setSelectedRobot(robot);
    disclosure.onOpen();
  }

  function handleClose() {
    disclosure.onClose();
    setSelectedRobot(null);
  }

  return (
    <LStack>
      <LStack>
        {data.robots.map((robot) => (
          <RobotCard
            key={robot.id}
            robot={robot}
            onClick={() => handleRobotClick(robot)}
          />
        ))}
      </LStack>

      {selectedRobot && (
        <RobotConfigurationModal
          robotId={selectedRobot.id}
          isOpen={disclosure.isOpen}
          onClose={handleClose}
        />
      )}
    </LStack>
  );
}
