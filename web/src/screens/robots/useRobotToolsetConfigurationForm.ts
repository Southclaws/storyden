import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import {
  getRobotToolsetGetKey,
  getRobotToolsetsListKey,
  robotToolsetCreate,
  robotToolsetUpdate,
} from "@/api/openapi-client/robots";
import { RobotToolset } from "@/api/openapi-schema";

export type RobotToolsetFormProps = {
  toolset?: RobotToolset;
  onSave?: (toolset: RobotToolset) => void;
  onDelete?: () => Promise<void>;
};

const FormSchema = z.object({
  name: z.string().min(1, "Name is required"),
  description: z.string(),
  instruction: z.string(),
  tools: z.array(z.string()),
});

export type RobotToolsetForm = z.infer<typeof FormSchema>;

export function useRobotToolsetConfigurationForm({
  toolset,
  onSave,
}: RobotToolsetFormProps) {
  const { mutate } = useSWRConfig();
  const isCreating = !toolset;
  const isEditable = isCreating || toolset.editable;
  const form = useForm<RobotToolsetForm>({
    defaultValues: {
      name: toolset?.name ?? "",
      description: toolset?.description ?? "",
      instruction: toolset?.instruction ?? "",
      tools: toolset?.tools ?? [],
    },
    resolver: zodResolver(FormSchema),
  });

  const handleSave = form.handleSubmit(async (data) => {
    if (!isEditable) return;

    await handle(
      async () => {
        const saved = toolset
          ? await robotToolsetUpdate(toolset.id, data)
          : await robotToolsetCreate(data);
        onSave?.(saved);
      },
      {
        promiseToast: {
          loading: isCreating ? "Creating Toolset..." : "Saving Toolset...",
          success: isCreating ? "Toolset created" : "Toolset saved",
        },
        async cleanup() {
          await mutate(getRobotToolsetsListKey());
          if (toolset) {
            await mutate(getRobotToolsetGetKey(toolset.id));
          }
        },
      },
    );
  });

  return { form, isCreating, isEditable, handleSave };
}
