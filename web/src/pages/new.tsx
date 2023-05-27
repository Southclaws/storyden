import { useRouter } from "next/router";
import { z } from "zod";
import dynamic from "next/dynamic";

// NOTE: Render client side, probably unnecessary but quick fix to react errors.
const ComposeScreen = dynamic(
  () =>
    import("../screens/compose/ComposeScreen").then((mod) => mod.ComposeScreen),
  {
    ssr: false,
  }
);

export const DraftQuerySchema = z.object({
  id: z.string().optional(),
});
export type DraftQuery = z.infer<typeof DraftQuerySchema>;

export default function Page() {
  const { query } = useRouter();
  const { id } = DraftQuerySchema.parse(query);

  return <ComposeScreen editing={id} />;
}
