"use client";

import { useCopyToClipboard } from "@uidotdev/usehooks";

import { PostReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { WEB_ADDRESS } from "src/config";
import { isShareEnabled } from "src/utils/client";

export type Props = {
  thread: PostReference;
  onDelete?: () => void;
};

export function useFeedItemMenu(props: Props) {
  const account = useSession();
  const permalink = getPermalinkForThread(props.thread.slug);
  const [, copyToClipboard] = useCopyToClipboard();

  const shareEnabled = isShareEnabled();
  const deleteEnabled =
    (account?.admin || account?.id === props.thread.author.id) &&
    props.onDelete;

  async function share() {
    await navigator.share({
      title: `A post by ${props.thread.author.name}`,
      url: `#${props.thread.id}`,
      text: props.thread.description,
    });
  }

  function handleSelect({ value }: { value: string }) {
    switch (value) {
      case "copy-link":
        copyToClipboard(permalink);
        return;

      case "share":
        share();
        return;

      case "delete":
        props.onDelete?.();
        return;

      default:
        throw new Error("unknown handler");
    }
  }

  return {
    handleSelect,
    shareEnabled,
    deleteEnabled,
  };
}

function getPermalinkForThread(slug: string) {
  return `${WEB_ADDRESS}/t/${slug}`;
}
