import { zodResolver } from "@hookform/resolvers/zod";
import { ReactNode } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { reportCreate } from "@/api/openapi-client/reports";
import { DatagraphItemKind } from "@/api/openapi-schema";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { UseDisclosureProps } from "@/utils/useDisclosure";

import { Input } from "../ui/input";

const FormSchema = z.object({
  comment: z.string().max(1000).optional(),
});
type Form = z.infer<typeof FormSchema>;

export type ReportModalProps = UseDisclosureProps & {
  title: string;
  description: ReactNode;
  subject?: ReactNode;
  targetId: string;
  targetKind: DatagraphItemKind;
  submitLabel?: string;
  successMessage?: string;
  loadingMessage?: string;
};

export function ReportModal({
  title,
  description,
  subject,
  targetId,
  targetKind,
  submitLabel = "Submit report",
  successMessage = "Report submitted",
  loadingMessage = "Sending report...",
  ...disclosure
}: ReportModalProps) {
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      comment: "",
    },
  });

  const handleSubmit = form.handleSubmit(async ({ comment }) => {
    await handle(
      async () => {
        await reportCreate({
          target_id: targetId,
          target_kind: targetKind,
          comment: comment?.trim() || undefined,
        });
      },
      {
        promiseToast: {
          loading: loadingMessage,
          success: successMessage,
        },
      },
    );

    form.reset();
    disclosure.onClose?.();
  });

  function handleCancel() {
    form.reset();
    disclosure.onClose?.();
  }

  return (
    <ModalDrawer title={title} {...disclosure}>
      <styled.form
        as="form"
        gap="4"
        alignItems="stretch"
        onSubmit={handleSubmit}
        className={lstack()}
      >
        <styled.div color="text.subtle">{description}</styled.div>

        {subject && (
          <styled.div
            borderWidth="thin"
            borderRadius="md"
            borderColor="border.muted"
            background="background.inset"
            padding="3"
          >
            {subject}
          </styled.div>
        )}

        <FormControl>
          <FormLabel>Additional details</FormLabel>
          <Input
            {...form.register("comment")}
            placeholder="Optional context for moderators"
            maxLength={1000}
            resize="vertical"
          />
          <FormHelperText>
            Optional additional context to help moderators resolve this matter.
          </FormHelperText>
          <FormErrorText>
            {form.formState.errors.comment?.message}
          </FormErrorText>
        </FormControl>

        <WStack gap="2">
          <Button type="button" variant="outline" onClick={handleCancel}>
            Cancel
          </Button>
          <Button
            type="submit"
            colorPalette="red"
            loading={form.formState.isSubmitting}
          >
            {submitLabel}
          </Button>
        </WStack>
      </styled.form>
    </ModalDrawer>
  );
}
