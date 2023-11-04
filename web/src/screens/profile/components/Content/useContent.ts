import { useCollectionList } from "src/api/openapi/collections";
import { usePostSearch } from "src/api/openapi/posts";
import { PublicProfile } from "src/api/openapi/schemas";
import { useFeed } from "src/components/feed/useFeed";

export function useContent(props: PublicProfile) {
  const threads = useFeed({ author: props.handle });
  const posts = usePostSearch({ author: props.handle, kind: ["post"] });
  const collections = useCollectionList();

  if (!threads.data) return { ready: false as const, error: threads.error };
  if (!posts.data) return { ready: false as const, error: posts.error };
  if (!collections.data) return { ready: false as const, error: posts.error };

  return {
    ready: true as const,
    data: {
      threads: threads.data.threads,
      posts: posts.data.results,
      collections: collections.data.collections,
    },
    handlers: threads.handlers,
  };
}
