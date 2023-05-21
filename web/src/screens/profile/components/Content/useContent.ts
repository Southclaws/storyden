import { APIError, PublicProfile, ThreadList } from "src/api/openapi/schemas";
import { useThreadList } from "src/api/openapi/threads";

type ContentResponse =
  | { ready: false; error: void | APIError }
  | {
      ready: true;
      data: ThreadList;
    };

export function useContent(props: PublicProfile): ContentResponse {
  const threads = useThreadList({ author: props.handle });

  if (!threads.data) return { ready: false, error: threads.error };

  return { ready: true, data: threads.data.threads };
}
