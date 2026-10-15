import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { handle } from "@/api/client";
import { reportCreate } from "@/api/openapi-client/reports";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import * as z from "@/lib/validation";
import { WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Input } from "../ui/input";

import type { ReportModalProps } from "./ReportModal";

const FormSchema = z.object({
  comment: z.optional(z.string().check(z.maxLength(1000))),
});
type Form = z.infer<typeof FormSchema>;

type Props = Omit<ReportModalProps, "title"> & {
  initialComment: string;
  onCommentChange: (comment: string) => void;
};

export function ReportForm({
  initialComment,
  onCommentChange,
  description,
  subject,
  targetId,
  targetKind,
  submitLabel = "Submit report",
  successMessage = "Report submitted",
  loadingMessage = "Sending report...",
  ...disclosure
}: Props) {
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      comment: initialComment,
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

    onCommentChange("");
    form.reset({ comment: "" });
    disclosure.onClose?.();
  });

  function handleCancel() {
    onCommentChange("");
    form.reset({ comment: "" });
    disclosure.onClose?.();
  }

  return (
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
          {...form.register("comment", {
            onChange: (event) => {
              onCommentChange(event.target.value);
            },
          })}
          placeholder="Optional context for moderators"
          maxLength={1000}
          resize="vertical"
        />
        <FormHelperText>
          Optional additional context to help moderators resolve this matter.
        </FormHelperText>
        <FormErrorText>{form.formState.errors.comment?.message}</FormErrorText>
      </FormControl>

      <WStack gap="2">
        <Button type="button" variant="outline" onClick={handleCancel}>
          Cancel
        </Button>
        <Button
          type="submit"
          variant="solid"
          loading={form.formState.isSubmitting}
        >
          {submitLabel}
        </Button>
      </WStack>
    </styled.form>
  );
}
