"use client";

import { useRouter } from "next/navigation";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  getRobotToolsetGetKey,
  getRobotToolsetsListKey,
  robotToolsetDelete,
  useRobotToolsetGet,
} from "@/api/openapi-client/robots";
import { RobotToolsetConfigurationForm } from "@/components/robots/RobotToolsetConfiguration/RobotToolsetConfigurationForm";
import { BackAction } from "@/components/site/Action/Back";
import { UnreadyBanner } from "@/components/site/Unready";
import { PageHeading } from "@/components/ui/page-heading";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

export function RobotToolsetConfigurationScreen({
  toolsetId,
}: {
  toolsetId: string;
}) {
  const { data, error } = useRobotToolsetGet(toolsetId);
  const { mutate } = useSWRConfig();
  const router = useRouter();

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  async function handleDelete() {
    await handle(() => robotToolsetDelete(toolsetId), {
      promiseToast: {
        loading: "Deleting Toolset...",
        success: "Toolset deleted",
      },
      async cleanup() {
        await mutate(getRobotToolsetsListKey());
        await mutate(getRobotToolsetGetKey(toolsetId), undefined, {
          revalidate: false,
        });
      },
    });
    router.push("/robots/toolsets");
  }

  return (
    <LStack gap="4" w="full">
      <WStack>
        <HStack gap="2">
          <BackAction fallbackHref="/robots/toolsets" />
          <PageHeading lineClamp="1">{data.name}</PageHeading>
        </HStack>
      </WStack>
      <RobotToolsetConfigurationForm
        toolset={data}
        onDelete={data.editable ? handleDelete : undefined}
      />
    </LStack>
  );
}
