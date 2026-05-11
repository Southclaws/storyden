import { Unready } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { RadioGroupField } from "@/components/ui/form/RadioGroupField";
import { Heading } from "@/components/ui/heading";
import { useI18n } from "@/i18n/provider";
import { CardBox, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import {
  Props,
  useMemberInterfaceSettings,
} from "./useMemberInterfaceSettings";

export function MemberInterfaceSettings(props: Props) {
  const result = useMemberInterfaceSettings(props);
  const { t } = useI18n();

  if (!result.ready) {
    return <Unready />;
  }

  const { control, formState, onSubmit } = result;

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
          <FormLabel>{t("Text editor style")}</FormLabel>
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
              "Choose your preferred editor style for composing threads, replies and pages.",
            )}
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("Sidebar default state")}</FormLabel>
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
              "Choose your preferred default state for the sidebar when you visit the site.",
            )}
          </FormHelperText>
        </FormControl>
      </CardBox>
    </styled.form>
  );
}
