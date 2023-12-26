"use client";

import { useCopyToClipboard } from "@uidotdev/usehooks";

import { ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { WEB_ADDRESS } from "src/config";
import { isShareEnabled } from "src/utils/client";

export type Props = {
  thread: ThreadReference;
  onDelete: () => void;
};

export function useFeedItemMenu(props: Props) {
  const account = useSession();
  const permalink = getPermalinkForThread(props.thread.slug);
  const [, copyToClipboard] = useCopyToClipboard();

  const shareEnabled = isShareEnabled();
  const deleteEnabled =
    account?.admin || account?.id === props.thread.author.id;

  async function onCopyLink() {
    copyToClipboard(permalink);
  }

  async function onShare() {
    await navigator.share({
      title: `A post by ${props.thread.author.name}`,
      url: `#${props.thread.id}`,
      text: props.thread.short,
    });
  }

  return {
    onCopyLink,
    shareEnabled,
    onShare,
    deleteEnabled,
    onDelete: props.onDelete,
  };
}

function getPermalinkForThread(slug: string) {
  return `${WEB_ADDRESS}/t/${slug}`;
}
