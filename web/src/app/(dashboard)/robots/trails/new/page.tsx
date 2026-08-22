"use client";

import { BackAction } from "@/components/site/Action/Back";
import { PageHeader } from "@/components/ui/page-header";
import { TrailEditor } from "@/screens/trails/TrailEditor";
import { LStack } from "@/styled-system/jsx";

export default function Page() {
  return (
    <LStack gap="4">
      <PageHeader
        title="New Trail"
        back={<BackAction fallbackHref="/robots/trails" />}
      />

      <TrailEditor />
    </LStack>
  );
}
