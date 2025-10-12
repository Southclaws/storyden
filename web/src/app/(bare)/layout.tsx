import Image from "next/image";
import { PropsWithChildren } from "react";

import { Fullpage } from "src/layouts/Fullpage";

import { BackAction } from "@/components/site/Action/Back";
import { HomeAnchor } from "@/components/site/Navigation/Anchors/Home";
import { getSettings } from "@/lib/settings/settings-server";
import { css } from "@/styled-system/css";
import { CardBox, VStack, WStack, styled } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";
import { getIconURL } from "@/utils/icon";

export default async function Layout({ children }: PropsWithChildren) {
  const settings = await getSettings();

  const siteName = settings.title;

  return (
    <Fullpage>
      <VStack minH="dvh" py="24">
        <CardBox className={vstack()} maxW="sm" gap="4" p="4">
          <WStack>
            <BackAction />

            <HomeAnchor hideLabel />
          </WStack>
          <VStack>
            <Image
              className={css({ width: "24", borderRadius: "md" })}
              src={getIconURL("512x512")}
              width={512}
              height={512}
              alt={`The ${siteName} logo`}
            />

            <styled.h1 fontWeight="bold" fontSize="lg">
              {siteName}
            </styled.h1>
          </VStack>

          {children}
        </CardBox>
      </VStack>
    </Fullpage>
  );
}
