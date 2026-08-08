"use client";

import { useNodeDraftList, useNodeList } from "@/api/openapi-client/nodes";
import { Visibility } from "@/api/openapi-schema";
import { QueueNodeList } from "@/components/queue/QueueNodeList";
import { QueueVersionList } from "@/components/queue/QueueVersionList";
import { Unready } from "@/components/site/Unready";
import { PageHeading } from "@/components/ui/page-heading";
import { SectionHeading } from "@/components/ui/section-heading";
import { LStack } from "@/styled-system/jsx";

export function QueueScreen() {
  const { data: submissions, error: submissionError } = useNodeList({
    visibility: [Visibility.review],
    format: "flat",
  });
  const { data: drafts, error: draftsError } = useNodeDraftList();

  if (!submissions || !drafts) {
    return <Unready error={submissionError ?? draftsError} />;
  }

  return (
    <LStack gap="8">
      <LStack>
        <PageHeading>Submission queue</PageHeading>

        <QueueNodeList nodes={submissions.nodes} />
      </LStack>

      <LStack>
        <SectionHeading>Page edits for review</SectionHeading>

        <QueueVersionList drafts={drafts.drafts} />
      </LStack>
    </LStack>
  );
}
