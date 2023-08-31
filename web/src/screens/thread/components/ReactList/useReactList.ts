import { NativeRenderer, createPicker } from "picmo";
import { mutate } from "swr";

import { postReactAdd } from "src/api/openapi/posts";
import { PostProps } from "src/api/openapi/schemas";
import { getThreadGetKey } from "src/api/openapi/threads";
import { useSession } from "src/auth";

export const emojiPickerContainerID = `react-emoji-select`;

export type Props = PostProps & {
  slug?: string;
};

export function useReactList(props: Props) {
  const account = useSession();
  const authenticated = !!account;

  async function onSelect(event: { emoji: string }) {
    await postReactAdd(props.id, { emoji: event.emoji });
    props.slug && mutate(getThreadGetKey(props.slug));
  }

  async function onOpen() {
    const rootElement = document.querySelector(
      `#${emojiPickerContainerID}-${props.id}`
    ) as HTMLElement;

    if (!rootElement) {
      throw new Error("cannot find emoji picker container");
    }

    const picker = createPicker({
      rootElement,
      renderer: new NativeRenderer(),
      //   emojiSize: "1.8rem",
    });

    picker.addEventListener("emoji:select", onSelect);
  }

  return { onOpen, authenticated };
}
