import { useState } from "react";

import { Robot } from "@/api/openapi-schema";
import { TOOL_NAMES } from "@/api/robots";
import { FormControl } from "@/components/ui/FormControl";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { FormLabel } from "@/components/ui/FormLabel";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/MultiSelectPicker";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import {
  Props,
  useRobotConfigurationForm,
} from "../../../screens/admin/RobotsSettingsScreen/useRobotConfigurationForm";

const mapToolToPickerItem = (name) => ({
  label: name,
  value: name,
});

const TOOL_OPTIONS: MultiSelectPickerItem[] =
  TOOL_NAMES.map(mapToolToPickerItem);

export function RobotConfigurationForm(props: Props) {
  const {
    form,
    handlers: { handleSave },
  } = useRobotConfigurationForm(props);

  const [tools, setTools] = useState<MultiSelectPickerItem[]>(
    props.robot.tools.map(mapToolToPickerItem),
  );
  const [toolQuery, setToolQuery] = useState("");

  const filteredToolOptions = toolQuery
    ? TOOL_OPTIONS.filter((opt) =>
        opt.label.toLowerCase().includes(toolQuery.toLowerCase()),
      )
    : TOOL_OPTIONS;

  async function handleToolsChange(items: MultiSelectPickerItem[]) {
    setTools(items);
  }

  return (
    <styled.form
      className={lstack()}
      h="full"
      justifyContent="space-between"
      onSubmit={handleSave}
    >
      <LStack gap="4" overflowY="auto" px="0.5" pb="1">
        <FormControl>
          <FormLabel>Name</FormLabel>
          <Input {...form.register("name")} placeholder="Robot name" />
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Description</FormLabel>
          <Input
            {...form.register("description")}
            placeholder="What this robot does"
          />
          <FormErrorText>
            {form.formState.errors.description?.message}
          </FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Playbook</FormLabel>
          <styled.textarea
            {...form.register("playbook")}
            placeholder="Instructions for the robot..."
            rows={12}
            w="full"
            minH="48"
            p="3"
            fontFamily="mono"
            fontSize="sm"
            lineHeight="relaxed"
            borderWidth="thin"
            borderStyle="solid"
            borderColor="border.default"
            borderRadius="l2"
            resize="vertical"
            _focus={{
              outline: "none",
              borderColor: "border.accent",
            }}
          />
          <FormErrorText>
            {form.formState.errors.playbook?.message}
          </FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Tools</FormLabel>
          <MultiSelectPicker
            value={tools}
            onChange={handleToolsChange}
            onQuery={setToolQuery}
            queryResults={filteredToolOptions}
            inputPlaceholder="Select tools..."
            size="sm"
          />
        </FormControl>
      </LStack>

      <WStack>
        <Button flexGrow="1" disabled={!form.formState.isDirty}>
          Save
        </Button>
      </WStack>
    </styled.form>
  );
}
