import { usePostSearch } from "src/api/openapi/posts";
import {
  APIError,
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
      };
    };

export function useContent(props: PublicProfile): ContentResponse {
  const threads = useThreadList({ author: props.handle });
  const posts = usePostSearch({ author: props.handle, kind: ["post"] });

  if (!threads.data) return { ready: false, error: threads.error };
  if (!posts.data) return { ready: false, error: posts.error };

  return {
    ready: true,
    data: {
      threads: threads.data.threads,
      posts: posts.data.results,
    },
  };
}
