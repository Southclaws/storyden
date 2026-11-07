import { Unready } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form-control";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { PageHeading } from "@/components/ui/page-heading";
import { RadioGroupField } from "@/components/ui/radio-group";
import { Text } from "@/components/ui/text";
import { LStack, WStack, styled } from "@/styled-system/jsx";

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
      <LStack gap="1">
        <WStack>
          <PageHeading>Interface settings</PageHeading>
          <Button type="submit" loading={formState.isSubmitting}>
            Save
          </Button>
        </WStack>
        <Text variant="supporting">
          Customise how you view and experience this site.
        </Text>
      </LStack>

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
          Choose your preferred editor style for composing threads, replies and
          pages.
        </FormHelperText>
      </FormControl>
    </styled.form>
  );
}
