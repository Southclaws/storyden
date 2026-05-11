import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { useI18n } from "@/i18n/provider";
import { HStack, VStack, styled } from "@/styled-system/jsx";
import { relativeTimestamp } from "@/utils/date";

import { DeleteDeviceTrigger } from "./DeleteDevice/DeleteDeviceTrigger";
import { Props, useDevices } from "./useDevices";

export function Devices(props: Props) {
  const { handleDeviceRegister } = useDevices();
  const { t, locale } = useI18n();

  return (
    <VStack w="full" alignItems="start">
      <Heading size="sm">{t("Devices")}</Heading>

      <p>
        {t(
          "You can use certain support devices with biometric authentication. You can add as many as you want here. It's recommended that you also add at least one extra device or other authentication method in case you lose your device.",
        )}
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
              {t("Device ID")}:{" "}
              <styled.code title={v.identifier}> {v.identifier}</styled.code>
            </styled.p>

            <HStack justify="space-between">
              <styled.p>
                {t("Created")}{" "}
                <time>{relativeTimestamp(v.created_at, locale)}</time>
              </styled.p>

              <DeleteDeviceTrigger id={v.id} />
            </HStack>
          </styled.li>
        ))}
      </styled.ul>

      <Button variant="subtle" onClick={handleDeviceRegister}>
        {t("Register this device")}
      </Button>
    </VStack>
  );
}
