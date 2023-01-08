import { useRouter } from "next/router";
import { useThreadsGet } from "src/api/openapi/threads";
import { Unready } from "src/components/Unready";
import { ThreadView } from "./components/ThreadView.tsx/ThreadView";

export function ThreadScreen() {
  const router = useRouter();
  const slug = router.query["slug"] as string;

  const thread = useThreadsGet(slug);

  if (!thread.data) return <Unready {...thread.error} />;

  return <ThreadView {...thread.data} />;
}
