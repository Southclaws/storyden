import { useLinkList } from "src/api/openapi/links";
import { LinkListResult } from "src/api/openapi/schemas";

export type Props = {
  query?: string;
  page?: number;
  links: LinkListResult;
};

export function useLinkIndexScreen(props: Props) {
  const { data, mutate, error } = useLinkList(
    {
      q: props.query,
      page: props.page?.toString(),
    },
    {
      swr: {
        fallbackData: props.links,
      },
    },
  );

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  return {
    ready: true as const,
    data,
    mutate,
  };
}
