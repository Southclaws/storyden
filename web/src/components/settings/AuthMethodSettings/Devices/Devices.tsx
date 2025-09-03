import { formatDistanceToNow } from "date-fns";

import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { HStack, VStack, styled } from "@/styled-system/jsx";

import { DeleteDeviceTrigger } from "./DeleteDevice/DeleteDeviceTrigger";
import { Props, useDevices } from "./useDevices";

export function Devices(props: Props) {
  const { handleDeviceRegister } = useDevices();

  return (
    <VStack w="full" alignItems="start">
      <Heading size="sm">Devices</Heading>

      <p>
        You can use certain support devices with biometric authentication. You
        can add as many as you want here. It&apos;s recommended that you also
        add at least one extra device or other authentication method in case you
        lose your device.
      </p>

      <styled.ul w="full" display="flex" flexDir="column" gap="2">
        {props.active.map((v) => (
          <styled.li
            key={v.id}
            display="flex"
            flexDir="column"
            borderColor="border.muted"
            borderWidth="thin"
            borderRadius="md"
            p="2"
            gap="2"
            minW="0"
          >
            <HStack justify="space-between">
              <Heading size="xs">{v.name}</Heading>
            </HStack>

            <styled.p
              minW="0"
              className="typography"
              whiteSpace="nowrap"
              textOverflow="ellipsis"
              overflow="hidden"
            >
              Device ID:{" "}
              <styled.code title={v.identifier}> {v.identifier}</styled.code>
            </styled.p>

            <HStack justify="space-between">
              <styled.p>
                Created{" "}
                <time>{formatDistanceToNow(new Date(v.created_at))}</time> ago
              </styled.p>

              <DeleteDeviceTrigger id={v.id} />
            </HStack>
          </styled.li>
        ))}
      </styled.ul>

      <Button variant="subtle" onClick={handleDeviceRegister}>
        Register this device
      </Button>
    </VStack>
  );
}
