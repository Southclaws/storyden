import Image from "next/image";

import { Heading } from "@/components/ui/heading";
import { css } from "@/styled-system/css";
import { HStack, LStack } from "@/styled-system/jsx";
import { getIconURL } from "@/utils/icon";
import { getInfo } from "@/utils/info";

export async function RootContextPane() {
  const info = await getInfo();
  const iconURL = getIconURL("512x512");

  return (
    <LStack>
      <HStack w="full" justify="space-between" alignItems="start">
        <Heading textWrap="wrap">{info.title}</Heading>

        <Image
          className={css({
            borderRadius: "md",
          })}
          alt="Icon"
          src={iconURL}
          width={32}
          height={32}
        />
      </HStack>

      <p>{info.description}</p>
    </LStack>
  );
}
