import { useClickAway } from "@uidotdev/usehooks";
import { mutate } from "swr";

import { postReactAdd } from "src/api/openapi/posts";
import { PostProps } from "src/api/openapi/schemas";
import { getThreadGetKey } from "src/api/openapi/threads";
import { useSession } from "src/auth";
import { useDisclosure } from "src/utils/useDisclosure";

export type Props = PostProps & {
  slug?: string;
};

type EmojiSelectEvent = {
  native: string;
};

export function useReactList(props: Props) {
  const account = useSession();
  const authenticated = !!account;

  const { onOpen, onClose, isOpen } = useDisclosure();
  const ref = useClickAway<HTMLDivElement>(() => {
    onClose();
  });

  async function handleSelect(event: EmojiSelectEvent) {
    await postReactAdd(props.id, { emoji: event.native });
    props.slug && mutate(getThreadGetKey(props.slug));

    onClose();
  }

  function handleTrigger() {
    if (!authenticated) {
      return;
    }

    // NOTE: Doesn't currently work to close the popover if it's open, because
    // by the time this handler is called, the outside click handler has already
    // closed the popover. But... tbh who cares! Very low impact issue to fix...
    if (!isOpen) onOpen();
    else onClose();
  }

  return {
    authenticated,
    isOpen,
    ref,
    handlers: {
      handleTrigger,
      handleSelect,
    },
  };
}
