import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import { getRobotsListKey, robotUpdate } from "@/api/openapi-client/robots";
import { Robot } from "@/api/openapi-schema";

export type Props = {
  robot: Robot;
  onSave?: () => void;
};

export const FormSchema = z.object({
  name: z.string().min(1, "Name is required"),
  description: z.string(),
  playbook: z.string(),
  tools: z.array(z.string()),
});
export type Form = z.infer<typeof FormSchema>;

export function useRobotConfigurationForm({ robot, onSave }: Props) {
  const { mutate } = useSWRConfig();
  const form = useForm<Form>({
    defaultValues: {
      name: robot.name,
      description: robot.description,
      playbook: robot.playbook,
      tools: robot.tools,
    },
    resolver: zodResolver(FormSchema),
  });

  const revalidate = async () => {
    await mutate(getRobotsListKey());
  };

  const handleSave = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await robotUpdate(robot.id, data);
        onSave?.();
      },
      {
        promiseToast: {
          loading: "Saving robot...",
          success: "Robot saved",
        },
        async cleanup() {
          await revalidate();
        },
      },
    );
  });

  return {
    form,
    handlers: {
      handleSave,
    },
  };
}
