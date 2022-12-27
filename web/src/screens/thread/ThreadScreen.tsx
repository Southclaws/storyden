import { useRouter } from "next/router";
import { useThreadsGet } from "src/api/openapi/threads";
import { Unready } from "src/components/Unready";
import { Thread } from "./components/Thread";

export function ThreadScreen() {
  const router = useRouter();
  const slug = router.query["t"] as string;

  const thread = useThreadsGet(slug);

  if (!thread.data) return <Unready {...thread.error} />;

  return <Thread {...thread.data} />;
}
