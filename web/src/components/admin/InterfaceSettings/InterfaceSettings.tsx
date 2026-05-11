import { Controller } from "react-hook-form";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { NumberInputField } from "@/components/ui/form/NumberInputField";
import { RadioGroupField } from "@/components/ui/form/RadioGroupField";
import { Heading } from "@/components/ui/heading";
import { useI18n } from "@/i18n/provider";
import { CardBox, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Props, useInterfaceSettings } from "./useInterfaceSettings";

export function InterfaceSettingsForm(props: Props) {
  const { t } = useI18n();
  const { control, signaturesEnabled, formState, onSubmit } =
    useInterfaceSettings(props);

  return (
    <styled.form
      width="full"
      display="flex"
      flexDirection="column"
      gap="4"
      onSubmit={onSubmit}
    >
      <CardBox className={lstack()}>
        <WStack>
          <Heading size="md">{t("Interface settings")}</Heading>
          <Button type="submit" loading={formState.isSubmitting}>
            {t("Save")}
          </Button>
        </WStack>

        <FormControl>
          <FormLabel>{t("Default editor")}</FormLabel>
          <RadioGroupField
            control={control}
            name="editorMode"
            items={[
              { label: t("Rich text"), value: "richtext" },
              { label: t("Markdown"), value: "markdown" },
            ]}
          />
          <FormHelperText>
            {t(
              "Choose the default editor for composing threads, replies and pages.",
            )}
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("Default sidebar state")}</FormLabel>
          <RadioGroupField
            control={control}
            name="sidebarDefaultState"
            items={[
              { label: t("Open"), value: "open" },
              { label: t("Closed"), value: "closed" },
            ]}
          />
          <FormHelperText>
            {t(
              "Choose the default state for the sidebar when members first visit or when they haven't set a preference.",
            )}
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("Signatures")}</FormLabel>
          <Controller
            control={control}
            name="signaturesEnabled"
            render={({ field }) => (
              <Checkbox
                size="sm"
                checked={!!field.value}
                onCheckedChange={({ checked }) => {
                  field.onChange(checked === true);
                }}
              >
                {t("Enable member signatures")}
              </Checkbox>
            )}
          />
          <FormHelperText>
            {t(
              "When disabled, signatures are hidden under posts and on profiles.",
            )}
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("Signature max height (px)")}</FormLabel>
          <NumberInputField
            control={control}
            name="signatureMaxHeight"
            min={32}
            max={2000}
            disabled={!signaturesEnabled}
          />
          <FormHelperText>
            {t("Limits how tall member signatures can appear below posts.")}
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("Signature max characters")}</FormLabel>
          <NumberInputField
            control={control}
            name="signatureMaxChars"
            min={1}
            max={10000}
            disabled={!signaturesEnabled}
          />
          <FormHelperText>
            {t("Visible characters, not including HTML tags.")}
          </FormHelperText>
        </FormControl>
      </CardBox>
    </styled.form>
  );
}
