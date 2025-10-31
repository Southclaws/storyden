"use client";

import { useRouter } from "next/navigation";

import { Account, Collection } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { useSession } from "@/auth";
import { useConfirmation } from "@/components/site/useConfirmation";
import { useCollectionMutations } from "@/lib/collection/mutation";
import {
  canDeleteCollection,
  canEditCollection,
} from "@/lib/collection/permissions";
import { useFeedMutations } from "@/lib/feed/mutation";
import { useShare } from "@/utils/client";
import { useCopyToClipboard } from "@/utils/useCopyToClipboard";

export type Props = {
  session?: Account;
  collection: Collection;
};

export function useCollectionMenu({ session, collection }: Props) {
  const router = useRouter();
  const account = useSession();
  const [, copyToClipboard] = useCopyToClipboard();

  const { deleteCollection, revalidate } = useCollectionMutations(session);

  const isSharingEnabled = useShare();
  const isEditingEnabled = canEditCollection(collection, account);
  const isDeletingEnabled = canDeleteCollection(collection, account);

  const permalink = `/c/${collection.slug}`;

  async function handleCopyLink() {
    copyToClipboard(permalink);
  }

  async function handleShare() {
    await navigator.share({
      title: `A collection by ${collection.owner.name}`,
      url: permalink,
      text: collection.description,
    });
  }

  async function handleDelete() {
    await handle(
      async () => {
        await deleteCollection(collection.id);
      },
      {
        cleanup: async () => await revalidate(),
      },
    );
  }

  return {
    isSharingEnabled,
    isEditingEnabled,
    isDeletingEnabled,

    handlers: {
      handleCopyLink,
      handleShare,
      handleDelete,
    },
  };
}
