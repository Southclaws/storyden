import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form-control";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { NumberInputField } from "@/components/ui/number-input";
import { PageHeading } from "@/components/ui/page-heading";
import { Text } from "@/components/ui/text";
import { Flex, LStack, WStack, styled } from "@/styled-system/jsx";

import { ModerationWordListEditorField } from "./ModerationWordListEditor.field";
import { Form, Props, useModerationSettings } from "./useModerationSettings";

export function ModerationSettingsForm(props: Props) {
  const { control, formState, onSubmit } = useModerationSettings(props);

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
          <PageHeading>Moderation settings</PageHeading>
          <Button type="submit" loading={formState.isSubmitting}>
            Save
          </Button>
        </WStack>
        <Text variant="supporting">
          Set content limits and automatic word-based moderation rules.
        </Text>
      </LStack>

      <Flex
        flexDir={{
          base: "column",
          md: "row",
        }}
        gap="2"
      >
        <FormControl>
          <FormLabel>Thread content maximum length</FormLabel>
          <NumberInputField
            control={control}
            name="threadBodyMaxSize"
            ariaLabel="Thread content maximum length"
            scrubber={true}
            min={1}
            max={1_000_000}
            step={100}
          />
          <FormHelperText>
            The maximum amount of characters allowed in the body content of a
            thread.
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Reply content maximum length</FormLabel>
          <NumberInputField
            control={control}
            name="replyBodyMaxSize"
            ariaLabel="Reply content maximum length"
            scrubber={true}
            min={1}
            max={1_000_000}
            step={100}
          />
          <FormHelperText>
            The maximum amount of characters allowed in the body content of a
            reply.
          </FormHelperText>
        </FormControl>
      </Flex>

      <Flex
        flexDir={{
          base: "column",
          md: "row",
        }}
        gap="2"
      >
        <FormControl>
          <FormLabel>Word report list</FormLabel>

          <ModerationWordListEditorField
            control={control}
            name="wordReportList"
          />

          <FormHelperText>
            Words and phrases that will automatically report and hide posts that
            contain them.
          </FormHelperText>
        </FormControl>
        <FormControl>
          <FormLabel>Word block list</FormLabel>

          <ModerationWordListEditorField
            control={control}
            name="wordBlockList"
          />

          <FormHelperText>
            Words and phrases that will instantly reject posts without creating
            a report.
          </FormHelperText>
        </FormControl>
      </Flex>
    </styled.form>
  );
}
