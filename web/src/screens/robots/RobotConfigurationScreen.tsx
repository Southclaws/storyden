"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { mutateTransaction } from "@/api/mutate";
import {
  getRobotGetKey,
  getRobotsListKey,
  robotDelete,
  useRobotGet,
} from "@/api/openapi-client/robots";
import { RobotsListOKResponse } from "@/api/openapi-schema";
import { RobotConfigurationForm } from "@/components/robots/RobotConfiguration/RobotConfigurationForm";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { ArrowLeftIcon } from "@/components/ui/icons/Arrow";
import { LinkButton } from "@/components/ui/link-button";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

type Props = {
  robotId: string;
};

export function RobotConfigurationScreen({ robotId }: Props) {
  const { data, error } = useRobotGet(robotId);
  const { mutate } = useSWRConfig();
  const router = useRouter();

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  async function handleDelete() {
    await handle(
      async () => {
        await mutateTransaction(
          mutate,
          [
            {
              key: getRobotsListKey(),
              optimistic: (current: RobotsListOKResponse | undefined) =>
                current
                  ? {
                      ...current,
                      robots: current.robots.filter(
                        (robot) => robot.id !== robotId,
                      ),
                      results: Math.max(0, current.results - 1),
                    }
                  : current,
            },
          ],
          () => robotDelete(robotId),
          { revalidate: true },
        );
        await mutate(getRobotGetKey(robotId), undefined, {
          revalidate: false,
        });
        router.push("/robots");
      },
      {
        promiseToast: {
          loading: "Deleting robot...",
          success: "Robot deleted",
        },
      },
    );
  }

  return (
    <LStack gap="4" w="full">
      <WStack>
        <HStack gap="2">
          <Link href="/robots">
            <IconButton variant="ghost" size="xs">
              <ArrowLeftIcon />
            </IconButton>
          </Link>

          <Heading size="md" lineClamp="1">
            {data.name}
          </Heading>
        </HStack>

        <LinkButton variant="subtle" size="xs" href="/robots/chats">
          Chats
        </LinkButton>
      </WStack>

      <RobotConfigurationForm robot={data} onDelete={handleDelete} />
    </LStack>
  );
}
