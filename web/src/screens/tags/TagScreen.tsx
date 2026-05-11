"use client";

import { useTagGet } from "@/api/openapi-client/tags";
import { Tag, TagName } from "@/api/openapi-schema";
import { DatagraphItemCard } from "@/components/datagraph/DatagraphItemCard";
import { Unready } from "@/components/site/Unready";
import { TagBadge } from "@/components/tag/TagBadge";
import { Breadcrumbs } from "@/components/ui/Breadcrumbs";
import { useI18n } from "@/i18n/provider";
import { HStack, LStack } from "@/styled-system/jsx";

type Props = {
  slug: TagName;
  initialTag: Tag;
};

export function TagScreen(props: Props) {
  const { t } = useI18n();
  const { data, error } = useTagGet(props.slug, {
    swr: { fallbackData: props.initialTag },
  });
  if (!data) {
    return <Unready error={error} />;
  }

  const tag = data;

  return (
    <LStack>
      <LStack gap="1">
        <Breadcrumbs
          index={{
            label: t("Tags"),
            href: "/tags",
          }}
          crumbs={[
            {
              label: tag.name,
              href: `/tags/${tag.name}`,
            },
          ]}
        />

        <HStack gap="1">
          <p>{t("Threads and library pages tagged with")}</p>
          <TagBadge tag={tag} />
        </HStack>
      </LStack>

      {tag.items.map((item) => (
        <DatagraphItemCard key={item.ref.id} item={item} />
      ))}
    </LStack>
  );
}
