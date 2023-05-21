import { useRouter } from "next/router";
import { useProfileGet } from "src/api/openapi/profiles";
import { APIError, PublicProfile } from "src/api/openapi/schemas";
import { z } from "zod";

export const ParamSchema = z.object({
  handle: z.string(),
});
export type Param = z.infer<typeof ParamSchema>;

type ProfileScreen =
  | { ready: false; error: void | APIError }
  | {
      ready: true;
      data: PublicProfile;
    };

export function useProfileScreen(): ProfileScreen {
  const router = useRouter();

  const { handle } = ParamSchema.parse(router.query);

  const { data, error } = useProfileGet(handle);

  if (!data) {
    return { error, ready: false };
  }

  return { data, ready: true };
}
