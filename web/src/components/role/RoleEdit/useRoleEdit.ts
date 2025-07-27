import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import {
  getRoleListKey,
  roleDelete,
  roleUpdate,
} from "@/api/openapi-client/roles";
import { Role } from "@/api/openapi-schema";
import { PermissionSchema } from "@/lib/permission/permission";

export type Props = {
  role: Role;
  onSave?: () => void;
};

export const FormSchema = z.object({
  name: z.string(),
  colour: z.string(),
  permissions: z.array(PermissionSchema),
});
export type Form = z.infer<typeof FormSchema>;

export function useRoleEditScreen({ role, onSave }: Props) {
  const { mutate } = useSWRConfig();
  const form = useForm<Form>({
    defaultValues: role,
    resolver: zodResolver(FormSchema),
  });

  const revalidate = async () => {
    await mutate(getRoleListKey());
  };

  const handleSave = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await roleUpdate(role.id, data);
        onSave?.();
      },
      {
        promiseToast: {
          loading: "Saving...",
          success: "Saved!",
        },
        async cleanup() {
          await revalidate();
        },
      },
    );
  });

  async function handleDelete() {
    await handle(
      async () => {
        await roleDelete(role.id);
      },
      {
        promiseToast: {
          loading: "Deleting...",
          success: "Deleted!",
        },
        async cleanup() {
          await revalidate();
        },
      },
    );
  }

  async function handleReset() {
    await handle(
      async () => {
        await roleDelete(role.id);
      },
      {
        promiseToast: {
          loading: "Resetting...",
          success: "Reset to defaults!",
        },
        async cleanup() {
          await revalidate();
        },
      },
    );
  }

  return {
    form,
    handlers: {
      handleSave,
      handleDelete,
      handleReset,
    },
  };
}
