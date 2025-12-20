import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { RadioGroupField } from "@/components/ui/form/RadioGroupField";
import { Heading } from "@/components/ui/heading";
import { CardBox, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Props, useInterfaceSettings } from "./useInterfaceSettings";

export function InterfaceSettingsForm(props: Props) {
  const { control, formState, onSubmit } = useInterfaceSettings(props);

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
          <Heading size="md">Interface settings</Heading>
          <Button type="submit" loading={formState.isSubmitting}>
            Save
          </Button>
        </WStack>

        <FormControl>
          <FormLabel>Default editor</FormLabel>
          <RadioGroupField
            control={control}
            name="editorMode"
            items={[
              { label: "Rich text", value: "richtext" },
              { label: "Markdown", value: "markdown" },
            ]}
          />
          <FormHelperText>
            Choose the default editor for composing threads, replies and pages.
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Default sidebar state</FormLabel>
          <RadioGroupField
            control={control}
            name="sidebarDefaultState"
            items={[
              { label: "Open", value: "open" },
              { label: "Closed", value: "closed" },
            ]}
          />
          <FormHelperText>
            Choose the default state for the sidebar when members first visit or
            when they haven't set a preference.
          </FormHelperText>
        </FormControl>
      </CardBox>
    </styled.form>
  );
}
