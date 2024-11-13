"use client";

import { useLinkGet } from "@/api/openapi-client/links";
import { Link } from "@/api/openapi-schema";
import { LinkView } from "@/components/library/links/LinkView";
import { Unready } from "@/components/site/Unready";

export type Props = {
  initialLink: Link;
  slug: string;
};

export function LinkScreen(props: Props) {
  const { data, error } = useLinkGet(props.slug, {
    swr: { fallbackData: props.initialLink },
  });

  if (!data) return <Unready error={error} />;

  return <LinkView link={data} />;
}
