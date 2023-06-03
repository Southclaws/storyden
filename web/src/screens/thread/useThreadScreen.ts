import { useThreadGet } from "src/api/openapi/threads";
import { useThreadScreenState } from "./state";
import { useParams } from "next/navigation";
import { z } from "zod";

export const ParamSchema = z.object({
  slug: z.string(),
});
export type Param = z.infer<typeof ParamSchema>;

export function useThreadScreen() {
  const params = useParams();

  const { slug } = ParamSchema.parse(params);

  const { data, error } = useThreadGet(slug);

  const state = useThreadScreenState(data);

  return { state, data, error };
}
