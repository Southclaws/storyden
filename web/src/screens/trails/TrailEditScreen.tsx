"use client";

import { useRouter } from "next/navigation";

import { useTrailGet } from "@/api/openapi-client/trails";
import { BackAction } from "@/components/site/Action/Back";
import { Unready } from "@/components/site/Unready";
import { PageHeader } from "@/components/ui/page-header";
import { LStack } from "@/styled-system/jsx";

import { TrailEditor } from "./TrailEditor";

export function TrailEditScreen({ trailId }: { trailId: string }) {
  const router = useRouter();
  const trailQuery = useTrailGet(trailId);

  if (!trailQuery.data) return <Unready error={trailQuery.error} />;

  const trail = trailQuery.data;

  return (
    <LStack gap="6" alignItems="stretch">
      <PageHeader
        title="Edit Trail"
        back={<BackAction fallbackHref={`/robots/trails/${trail.id}`} />}
      />

      <TrailEditor
        trail={trail}
        onSaved={async (saved) => {
          await trailQuery.mutate(saved, { revalidate: false });
          router.push(`/robots/trails/${saved.id}`);
          router.refresh();
        }}
      />
    </LStack>
  );
}
