import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { type AdminSettings } from "@/lib/settings/settings";
import * as z from "@/lib/validation";
import { deriveError } from "@/utils/error";

export type Props = {
  settings: AdminSettings;
};

export const FormSchema = z.object({
  imageMaxSizeMiB: z
    .number()
    .check(
      z.minimum(1 / 1024 ** 2, "The size limit must allow at least one byte."),
      z.maximum(Math.floor(Number.MAX_SAFE_INTEGER / 1024 ** 2)),
    ),
  imageMaxMegapixels: z
    .number()
    .check(
      z.minimum(
        0.000001,
        "The resolution limit must allow at least one pixel.",
      ),
      z.maximum(Math.floor(Number.MAX_SAFE_INTEGER / 1_000_000)),
    ),
});
export type Form = z.infer<typeof FormSchema>;

export function useMediaSettings({ settings }: Props) {
  const { revalidate, updateSettings } = useSettingsMutation();
  const maxSize = settings.services?.media?.image_max_size;
  const maxPixels = settings.services?.media?.image_max_pixels;

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      imageMaxSizeMiB: maxSize === undefined ? undefined : maxSize / 1024 ** 2,
      imageMaxMegapixels:
        maxPixels === undefined ? undefined : maxPixels / 1_000_000,
    },
  });

  const onSubmit = form.handleSubmit(async (data) => {
    form.clearErrors("root");

    await handle(
      async () => {
        await updateSettings({
          services: {
            media: {
              image_max_size: Math.round(data.imageMaxSizeMiB * 1024 ** 2),
              image_max_pixels: Math.round(data.imageMaxMegapixels * 1_000_000),
            },
          },
        });

        form.reset(data);
        await revalidate();
        toast.success("Settings saved");
      },
      {
        errorToast: false,
        onError: async (error) => {
          form.setError("root", { message: deriveError(error) });
        },
      },
    );
  });

  return {
    clearErrors: form.clearErrors,
    control: form.control,
    formState: form.formState,
    onSubmit,
  };
}
