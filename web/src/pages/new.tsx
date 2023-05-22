import { useRouter } from "next/router";
import { ComposeScreen } from "src/screens/compose/ComposeScreen";
import { z } from "zod";

export const DraftQuerySchema = z.object({
  id: z.string().optional(),
});
export type DraftQuery = z.infer<typeof DraftQuerySchema>;

export default function Page() {
  const { query } = useRouter();
  const { id } = DraftQuerySchema.parse(query);

  console.log("EDITING", id);

  return <ComposeScreen editing={id} />;
}
