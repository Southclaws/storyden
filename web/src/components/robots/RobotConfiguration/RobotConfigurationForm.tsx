import { useState } from "react";

import { useRobotModelsList } from "@/api/openapi-client/robots";
import { TOOL_NAMES } from "@/api/robots";
import { RobotModelComboboxField } from "@/components/robots/RobotModelComboboxField";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { FormControl } from "@/components/ui/FormControl";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { FormLabel } from "@/components/ui/FormLabel";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/MultiSelectPicker";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { useDisclosure } from "@/utils/useDisclosure";

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
    isCreating,
    handlers: { handleSave },
  } = useRobotConfigurationForm(props);

  const [tools, setTools] = useState<MultiSelectPickerItem[]>(
    (props.robot?.tools ?? []).map(mapToolToPickerItem),
  );
  const [toolQuery, setToolQuery] = useState("");
  const { data: modelData, error: modelError } = useRobotModelsList();
  const models = modelData?.models ?? [];

  const filteredToolOptions = toolQuery
    ? TOOL_OPTIONS.filter((opt) =>
        opt.label.toLowerCase().includes(toolQuery.toLowerCase()),
      )
    : TOOL_OPTIONS;

  async function handleToolsChange(items: MultiSelectPickerItem[]) {
    setTools(items);
    form.setValue(
      "tools",
      items.map((item) => item.value),
      { shouldDirty: true, shouldValidate: true },
    );
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
          <FormLabel>Model</FormLabel>
          <RobotModelComboboxField
            control={form.control}
            name="model"
            models={models}
            placeholder={isCreating ? "Use default model" : "Select a model"}
            disabled={!modelData || models.length === 0}
          />
          {modelError ? (
            <FormErrorText>Failed to load robot models.</FormErrorText>
          ) : (
            <styled.p color="fg.muted" fontSize="sm">
              {isCreating
                ? "Leave unset to use the configured default model."
                : "Choose one of the enabled provider models."}
            </styled.p>
          )}
          <FormErrorText>{form.formState.errors.model?.message}</FormErrorText>
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

      <WStack justifyContent="space-between">
        {props.onDelete ? (
          <RobotDeleteButton onDelete={props.onDelete} />
        ) : (
          <styled.span />
        )}

        <Button
          size="sm"
          minW="32"
          disabled={!isCreating && !form.formState.isDirty}
        >
          {isCreating ? "Create" : "Save"}
        </Button>
      </WStack>
    </styled.form>
  );
}

function RobotDeleteButton({ onDelete }: { onDelete: () => Promise<void> }) {
  const disclosure = useDisclosure();
  const [isDeleting, setIsDeleting] = useState(false);

  async function handleDelete() {
    setIsDeleting(true);

    try {
      await onDelete();
      disclosure.onClose();
    } finally {
      setIsDeleting(false);
    }
  }

  return (
    <>
      <Button
        type="button"
        size="sm"
        minW="32"
        variant="ghost"
        colorPalette="red"
        onClick={disclosure.onOpen}
      >
        Delete robot
      </Button>

      <ModalDrawer
        title="Delete robot"
        isOpen={disclosure.isOpen}
        onClose={disclosure.onClose}
      >
        <LStack gap="6">
          <LStack gap="2">
            <styled.p fontSize="sm">
              This will permanently delete this robot.
            </styled.p>
            <styled.p fontSize="sm" color="fg.muted">
              Existing robot chat sessions will remain, but this robot will no
              longer be available.
            </styled.p>
          </LStack>

          <HStack justifyContent="end" gap="3">
            <Button
              type="button"
              variant="ghost"
              onClick={disclosure.onClose}
              disabled={isDeleting}
            >
              Cancel
            </Button>
            <Button
              type="button"
              colorPalette="red"
              loading={isDeleting}
              onClick={handleDelete}
            >
              Delete robot
            </Button>
          </HStack>
        </LStack>
      </ModalDrawer>
    </>
  );
}
