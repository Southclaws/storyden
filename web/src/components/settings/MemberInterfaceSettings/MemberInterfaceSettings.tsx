import { Unready } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { RadioGroupField } from "@/components/ui/form/RadioGroupField";
import { Heading } from "@/components/ui/heading";
import { CardBox, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import {
  Props,
  useMemberInterfaceSettings,
} from "./useMemberInterfaceSettings";

export function MemberInterfaceSettings(props: Props) {
  const result = useMemberInterfaceSettings(props);

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
          <Heading size="md">Interface settings</Heading>
          <Button type="submit" loading={formState.isSubmitting}>
            Save
          </Button>
        </WStack>

        <FormControl>
          <FormLabel>Text editor style</FormLabel>
          <RadioGroupField
            control={control}
            name="editorMode"
            items={[
              { label: "Rich text", value: "richtext" },
              { label: "Markdown", value: "markdown" },
            ]}
          />
          <FormHelperText>
            Choose your preferred editor style for composing threads, replies
            and pages.
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Sidebar default state</FormLabel>
          <RadioGroupField
            control={control}
            name="sidebarDefaultState"
            items={[
              { label: "Open", value: "open" },
              { label: "Closed", value: "closed" },
            ]}
          />
          <FormHelperText>
            Choose your preferred default state for the sidebar when you visit
            the site.
          </FormHelperText>
        </FormControl>
      </CardBox>
    </styled.form>
  );
}
