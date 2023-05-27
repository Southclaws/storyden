import { useRouter } from "next/router";
import { Unready } from "src/components/Unready";
import { ProfileScreen } from "src/screens/profile/ProfileScreen";
import { z } from "zod";

export const ParamSchema = z.object({
  handle: z.string().optional(),
});
export type Param = z.infer<typeof ParamSchema>;

export default function Page() {
  const router = useRouter();
  const { handle } = ParamSchema.parse(router.query);

  if (!handle) return <Unready />;

  return <ProfileScreen handle={handle} />;
}
