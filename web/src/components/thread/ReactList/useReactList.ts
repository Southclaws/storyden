import { mutate } from "swr";

import { postReactAdd } from "src/api/openapi/posts";
import { PostProps } from "src/api/openapi/schemas";
import { getThreadGetKey } from "src/api/openapi/threads";
import { useSession } from "src/auth";

export type Props = PostProps & {
  slug?: string;
};

type EmojiSelectEvent = {
  native: string;
};

export function useReactList(props: Props) {
  const account = useSession();
  const authenticated = !!account;

  async function onSelect(event: EmojiSelectEvent) {
    await postReactAdd(props.id, { emoji: event.native });
    props.slug && mutate(getThreadGetKey(props.slug));
  }

  return {
    authenticated,
    handlers: {
      onSelect,
    },
  };
}
