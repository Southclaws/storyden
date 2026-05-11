import { ClipboardIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import * as Clipboard from "@/components/ui/clipboard";
import { DatePickerInputField } from "@/components/ui/form/DatePickerField";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { Input } from "@/components/ui/input";
import { getAPIAddress } from "@/config";
import { useI18n } from "@/i18n/provider";
import { LStack, WStack } from "@/styled-system/jsx";

import {
  Form,
  Props,
  useCreateAccessKeyScreen,
} from "./useCreateAccessKeyScreen";

export function CreateAccessKeyScreen({ onClose }: Props) {
  const { form, createdSecret, handleSubmit } = useCreateAccessKeyScreen({
    onClose,
  });
  const { t } = useI18n();

  if (createdSecret) {
    return (
      <LStack h="full" gap="8" justifyContent="space-between">
        <LStack>
          <Heading>{t("Access Key Created")}</Heading>
          <p>
            <strong>
              {t("This is the only time you'll see this access key.")}
            </strong>{" "}
            {t("Make sure to copy and store it securely.")}
          </p>

          <WStack>
            <Clipboard.Root w="full" value={createdSecret}>
              <Clipboard.Control>
                <Clipboard.Input asChild>
                  <Input />
                </Clipboard.Input>
                <Clipboard.Trigger asChild>
                  <IconButton variant="outline">
                    <Clipboard.Indicator copied={<CheckIcon />}>
                      <ClipboardIcon />
                    </Clipboard.Indicator>
                  </IconButton>
                </Clipboard.Trigger>
              </Clipboard.Control>
            </Clipboard.Root>
          </WStack>

          <LStack>
            <p>
              {t(
                "To use your new access key with this site, include it as an Authorization header with requests to the API.",
              )}
            </p>

            <Heading size="sm" color="fg.muted">
              {t("Header format:")}
            </Heading>

            <pre>Authorization: Bearer {createdSecret}</pre>

            <Heading size="sm" color="fg.muted">
              {t("API and MCP endpoints:")}
            </Heading>

            <pre>
              {getAPIAddress()}/api
              <br />
              {getAPIAddress()}/mcp
            </pre>
          </LStack>
        </LStack>

        <WStack>
          <Button onClick={onClose} style={{ width: "100%" }}>
            {t("Done")}
          </Button>
        </WStack>
      </LStack>
    );
  }

  return (
    <form onSubmit={handleSubmit}>
      <LStack gap="4">
        <div>
          <h3 style={{ marginBottom: "0.5rem", fontWeight: "600" }}>
            {t("Create Access Key")}
          </h3>
          <p style={{ fontSize: "0.875rem", color: "var(--colors-gray-600)" }}>
            {t(
              "Access keys allow you to authenticate API requests. They share the same permissions as your account.",
            )}
          </p>
        </div>

        <FormControl>
          <FormLabel>{t("Name")}</FormLabel>
          <Input
            {...form.register("name")}
            placeholder={t("e.g., Mobile App, CI/CD Pipeline")}
          />
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("Expiry Date (Optional)")}</FormLabel>
          <DatePickerInputField<Form>
            name="expires_at"
            control={form.control}
            // min={now("UTC")}
            // max={now("UTC").add({ years: 1 })}
          />
          <FormErrorText>
            {form.formState.errors.expires_at?.message}
          </FormErrorText>
          <p style={{ fontSize: "0.75rem", color: "var(--colors-gray-500)" }}>
            {t("Leave empty for no expiry")}
          </p>
        </FormControl>

        <WStack>
          <Button
            type="button"
            variant="outline"
            onClick={onClose}
            style={{ flex: 1 }}
          >
            {t("Cancel")}
          </Button>
          <Button
            type="submit"
            disabled={form.formState.isSubmitting}
            style={{ flex: 1 }}
          >
            {t("Create Key")}
          </Button>
        </WStack>
      </LStack>
    </form>
  );
}
