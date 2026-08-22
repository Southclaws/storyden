"use client";

import { Suspense } from "react";

import { Unready } from "@/components/site/Unready";
import { LinkButton } from "@/components/ui/link-button";
import { PageHeader } from "@/components/ui/page-header";
import { TrailListScreen } from "@/screens/trails/TrailListScreen";
import { LStack } from "@/styled-system/jsx";

export default function Page() {
  return (
    <LStack>
      <PageHeader
        title="Trails"
        description="Automate common tasks based on schedules or events."
        actions={<LinkButton href="/robots/trails/new">New Trail</LinkButton>}
      />
      <Suspense fallback={<Unready />}>
        <TrailListScreen />
      </Suspense>
    </LStack>
  );
}
