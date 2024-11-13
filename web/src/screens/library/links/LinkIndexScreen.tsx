"use client";

import { useLinkList } from "@/api/openapi-client/links";
import { LinkListResult } from "@/api/openapi-schema";
import { LinkIndexView } from "@/components/library/links/LinkIndexView/LinkIndexView";
import { Unready } from "@/components/site/Unready";

export type Props = {
  query?: string;
  page?: number;
  initialResult?: LinkListResult;
};

export function LinkIndexScreen(props: Props) {
  const { data, mutate, error } = useLinkList(
    {
      q: props.query,
      page: props.page?.toString(),
    },
    {
      swr: {
        fallbackData: props.initialResult,
      },
    },
  );

  if (!data) {
    return <Unready error={error} />;
  }

  return (
    <LinkIndexView
      links={data}
      mutate={mutate}
      query={props.query}
      page={props.page}
    />
  );
}
