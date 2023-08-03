import { useCollectionList } from "src/api/openapi/collections";
import { usePostSearch } from "src/api/openapi/posts";
import {
  APIError,
  Collection,
  PostProps,
  PublicProfile,
  ThreadReference,
} from "src/api/openapi/schemas";
import { useThreadList } from "src/api/openapi/threads";

type ContentResponse =
  | { ready: false; error: void | APIError }
  | {
      ready: true;
      data: {
        threads: ThreadReference[];
        posts: PostProps[];
        collections: Collection[];
      };
    };

export function useContent(props: PublicProfile): ContentResponse {
  const threads = useThreadList({ author: props.handle });
  const posts = usePostSearch({ author: props.handle, kind: ["post"] });
  const collections = useCollectionList();

  if (!threads.data) return { ready: false, error: threads.error };
  if (!posts.data) return { ready: false, error: posts.error };
  if (!collections.data) return { ready: false, error: posts.error };

  return {
    ready: true,
    data: {
      threads: threads.data.threads,
      posts: posts.data.results,
      collections: collections.data.collections,
    },
  };
}
