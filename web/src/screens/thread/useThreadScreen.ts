import { useRouter } from "next/router";
import { useThreadGet } from "src/api/openapi/threads";
import { useThreadScreenState } from "./state";

export function useThreadScreen() {
  const router = useRouter();
  const slug = router.query["slug"] as string;

  const { data, error } = useThreadGet(slug);

  const state = useThreadScreenState(data);

  return { state, data, error };
}
