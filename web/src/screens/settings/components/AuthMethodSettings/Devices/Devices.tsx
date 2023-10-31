import { Heading } from "@chakra-ui/react";

import { Button } from "src/theme/components/Button";

import { VStack, styled } from "@/styled-system/jsx";

import { Props, useDevices } from "./useDevices";

export function Devices(props: Props) {
  const { handleDeviceRegister } = useDevices();

  return (
    <VStack alignItems="start">
      <Heading size="sm">Devices</Heading>

      <p>
        You can use certain support devices with biometric authentication. You
        can add as many as you want here. It&apos;s recommended that you also
        add at least one extra device or other authentication method in case you
        lose your device.
      </p>

      <styled.ul>
        {props.active.map((v) => (
          <styled.li key={v.id}>{v.name}</styled.li>
        ))}
      </styled.ul>

      <Button onClick={handleDeviceRegister}>Register this device</Button>
    </VStack>
  );
}
