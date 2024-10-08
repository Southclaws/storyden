import { useCollectionList } from "src/api/openapi-client/collections";
import { PublicProfile } from "src/api/openapi-schema";

import { useThreadList } from "@/api/openapi-client/threads";

export type Props = {
  profile: PublicProfile;
};

export function useProfileContent({ profile }: Props) {
  const threads = useThreadList({ author: profile.handle });
  const collections = useCollectionList({ account_handle: profile.handle });

  if (!threads.data) {
    return { ready: false as const, error: threads.error };
  }
  if (!collections.data) {
    return { ready: false as const, error: collections.error };
  }

  return {
    ready: true as const,
    data: {
      threads: threads.data.threads,
      collections: collections.data.collections,
    },
  };
}
