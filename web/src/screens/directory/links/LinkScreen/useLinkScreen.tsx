import { useLinkList } from "src/api/openapi/links";
import { Link } from "src/api/openapi/schemas";

export type Props = {
  slug: string;
};

export function LinkScreen(props: Props) {
  const { data, error } = useLink(
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
