"use client";

import { useCopyToClipboard } from "@uidotdev/usehooks";

import { threadDelete } from "@/api/openapi-client/threads";
import { PostReference } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { WEB_ADDRESS } from "@/config";
import { useFeedMutations } from "@/lib/feed/mutation";
import { isShareEnabled } from "@/utils/client";

export type Props = {
  thread: PostReference;
};

export function useFeedItemMenu(props: Props) {
  const account = useSession();
  const permalink = getPermalinkForThread(props.thread.slug);
  const [, copyToClipboard] = useCopyToClipboard();

  const mutate = useFeedMutations();

  const shareEnabled = isShareEnabled();
  const deleteEnabled =
    account?.admin || account?.id === props.thread.author.id;

  async function share() {
    await navigator.share({
      title: `A post by ${props.thread.author.name}`,
      url: `#${props.thread.id}`,
      text: props.thread.description,
    });
  }

  async function handleDeleteThread() {
    await threadDelete(props.thread.id);
    mutate();
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
        handleDeleteThread();
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
