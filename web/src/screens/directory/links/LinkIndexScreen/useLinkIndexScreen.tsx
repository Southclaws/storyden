import { useLinkList } from "src/api/openapi/links";
import { LinkList } from "src/api/openapi/schemas";

export type Props = {
  query?: string;
  page?: number;
  links: LinkList;
};

export function useLinkIndexScreen(props: Props) {
  const { data, error } = useLinkList(
    {
      q: props.query,
      page: props.page?.toString(),
    },
    {
      swr: {
        fallbackData: { links: props.links },
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
    data: data.links,
  };
}
