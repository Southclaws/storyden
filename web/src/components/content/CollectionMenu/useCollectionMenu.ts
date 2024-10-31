import { useCollectionList } from "src/api/openapi-client/collections";
import { Account, Collection, PostReference } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { useCollectionItemMutations } from "@/lib/collection/mutation";

export type Props = {
  account: Account;
  thread: PostReference;
};

export function useCollectionMenu({ account, thread }: Props) {
  const { data, error } = useCollectionList({
    account_handle: account.handle,
    has_item: thread.id,
  });

  const { addPost, removePost, revalidate } =
    useCollectionItemMutations(account);

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const { collections } = data;

  const handleSelect = (collection: Collection) => async () => {
    await handle(
      async () => {
        const isAlreadySavedIn = collection?.has_queried_item;

        if (isAlreadySavedIn) {
          await removePost(collection.id, thread.id);
        } else {
          await addPost(collection, thread.id);
        }
      },
      { cleanup: async () => await revalidate() },
    );
  };

  return {
    ready: true as const,
    collections,
    handleSelect,
  };
}
