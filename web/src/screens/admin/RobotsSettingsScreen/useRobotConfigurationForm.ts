import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import {
  getRobotsListKey,
  robotCreate,
  robotUpdate,
} from "@/api/openapi-client/robots";
import { Robot, RobotCreateOKResponse } from "@/api/openapi-schema";

export type Props = {
  robot?: Robot;
  onSave?: (robot: RobotCreateOKResponse) => void;
  onDelete?: () => Promise<void>;
};

export const FormSchema = z.object({
  name: z.string().min(1, "Name is required"),
  description: z.string(),
  playbook: z.string(),
  model: z.string().optional(),
  tools: z.array(z.string()),
});
export type Form = z.infer<typeof FormSchema>;

export function useRobotConfigurationForm({ robot, onSave }: Props) {
  const { mutate } = useSWRConfig();
  const isCreating = !robot;
  const form = useForm<Form>({
    defaultValues: {
      name: robot?.name ?? "",
      description: robot?.description ?? "",
      playbook: robot?.playbook ?? "",
      model: robot?.model,
      tools: robot?.tools ?? [],
    },
    resolver: zodResolver(FormSchema),
  });

  const revalidate = async () => {
    await mutate(getRobotsListKey());
  };

  const handleSave = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        const payload = {
          ...data,
          model: data.model || undefined,
        };
        const saved = robot
          ? await robotUpdate(robot.id, payload)
          : await robotCreate(payload);
        onSave?.(saved);
      },
      {
        promiseToast: {
          loading: isCreating ? "Creating robot..." : "Saving robot...",
          success: isCreating ? "Robot created" : "Robot saved",
        },
        async cleanup() {
          await revalidate();
        },
      },
    );
  });

  return {
    form,
    isCreating,
    handlers: {
      handleSave,
    },
  };
}
