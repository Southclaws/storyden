"use client";

import { useNodeList } from "@/api/openapi-client/nodes";
import { Visibility } from "@/api/openapi-schema";
import { QueueNodeList } from "@/components/queue/QueueNodeList";
import { Unready } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { LStack } from "@/styled-system/jsx";

export function QueueScreen() {
  const { data, error } = useNodeList({
    visibility: [Visibility.review],
    format: "flat",
  });
  if (!data) {
    return <Unready error={error} />;
  }

  return (
    <LStack>
      <Heading>Submission queue</Heading>

      <QueueNodeList nodes={data.nodes} />
    </LStack>
  );
}
