"use client";

import { useCopyToClipboard } from "@uidotdev/usehooks";

import { PostReference } from "src/api/openapi-schema";
import { useSession } from "src/auth";
import { WEB_ADDRESS } from "src/config";
import { isShareEnabled } from "src/utils/client";

import { threadDelete } from "@/api/openapi-client/threads";

import { useFeedMutation } from "../../useFeed";

export type Props = {
  thread: PostReference;
};

export function useFeedItemMenu(props: Props) {
  const account = useSession();
  const permalink = getPermalinkForThread(props.thread.slug);
  const [, copyToClipboard] = useCopyToClipboard();

  const mutate = useFeedMutation();

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
