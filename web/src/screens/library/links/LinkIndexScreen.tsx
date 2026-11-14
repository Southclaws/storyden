"use client";

import { LinkListResult } from "@/api/openapi-schema";
import { LinkIndexView } from "@/components/library/links/LinkIndexView/LinkIndexView";
import { PageHeader } from "@/components/ui/page-header";
import { LStack } from "@/styled-system/jsx";

export type Props = {
  query?: string;
  page?: number;
  initialResult?: LinkListResult;
};

export function LinkIndexScreen(props: Props) {
  return (
    <LStack gap="4">
      <PageHeader
        title="Links"
        description="All links shared around your community."
      />
      <LinkIndexView {...props} />
    </LStack>
  );
}
