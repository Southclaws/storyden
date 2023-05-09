import { NativeRenderer, createPicker } from "picmo";
import { postReactAdd } from "src/api/openapi/posts";
import { Post } from "src/api/openapi/schemas";
import { getThreadGetKey } from "src/api/openapi/threads";
import { useAuthProvider } from "src/auth/useAuthProvider";
import { mutate } from "swr";

export const emojiPickerContainerID = `react-emoji-select`;

export type Props = Post & {
  slug: string;
};

export function useReactList(props: Props) {
  const { account } = useAuthProvider();
  const authenticated = !!account;

  async function onSelect(event: { emoji: string }) {
    await postReactAdd(props.id, { emoji: event.emoji });

    const key = getThreadGetKey(props.slug);

    console.log({ key });

    mutate(key);
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
    });

    picker.addEventListener("emoji:select", onSelect);
  }

  return { onOpen, authenticated };
}
