import { useState } from "react";

import {
  useRobotModelsList,
  useRobotToolsList,
  useRobotToolsetsList,
} from "@/api/openapi-client/robots";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Button } from "@/components/ui/button";
import { ComboboxField } from "@/components/ui/combobox";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormLabel } from "@/components/ui/form-label";
import { Input } from "@/components/ui/input";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/multi-select-picker";
import { Text } from "@/components/ui/text";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { useDisclosure } from "@/utils/useDisclosure";

import {
  Props,
  useRobotConfigurationForm,
} from "../../../screens/admin/RobotsSettingsScreen/useRobotConfigurationForm";

function collectToolsetTools(
  toolsets: { id: string; tools: string[] }[],
  selectedIDs: string[],
) {
  const selected = new Set(selectedIDs);
  const tools = new Set<string>();
  for (const toolset of toolsets) {
    if (!selected.has(toolset.id)) continue;
    for (const tool of toolset.tools) tools.add(tool);
  }
  return tools;
}

export function RobotConfigurationForm(props: Props) {
  const {
    form,
    isCreating,
    handlers: { handleSave },
  } = useRobotConfigurationForm(props);

  const [toolsetQuery, setToolsetQuery] = useState("");
  const [toolQuery, setToolQuery] = useState("");
  const { data: modelData, error: modelError } = useRobotModelsList();
  const { data: toolData, error: toolError } = useRobotToolsList();
  const { data: toolsetData, error: toolsetError } = useRobotToolsetsList();
  const models = modelData?.models ?? [];
  const availableTools = toolData?.tools ?? [];
  const availableToolsets = toolsetData?.toolsets ?? [];
  const selectedToolIDs = form.watch("tools");
  const selectedToolsetIDs = form.watch("toolsets");
  const selectedTools: MultiSelectPickerItem[] = selectedToolIDs.map((id) => {
    const tool = availableTools.find((candidate) => candidate.id === id);
    return { label: tool?.name ?? id, value: id };
  });
  const selectedToolsets: MultiSelectPickerItem[] = selectedToolsetIDs.map(
    (id) => {
      const toolset = availableToolsets.find(
        (candidate) => candidate.id === id,
      );
      return { label: toolset?.name ?? id, value: id };
    },
  );

  const normalizedToolsetQuery = toolsetQuery.trim().toLowerCase();
  const filteredToolsetOptions = availableToolsets.flatMap((toolset) => {
    if (
      normalizedToolsetQuery &&
      !toolset.name.toLowerCase().includes(normalizedToolsetQuery)
    ) {
      return [];
    }

    return [{ label: toolset.name, value: toolset.id }];
  });

  const selectedToolsetTools = collectToolsetTools(
    availableToolsets,
    selectedToolsetIDs,
  );
  const normalizedToolQuery = toolQuery.trim().toLowerCase();
  const filteredToolOptions = availableTools.flatMap((tool) => {
    if (!tool.available || selectedToolsetTools.has(tool.id)) {
      return [];
    }
    const label = tool.name ?? tool.id;
    if (
      normalizedToolQuery &&
      !label.toLowerCase().includes(normalizedToolQuery) &&
      !tool.id.toLowerCase().includes(normalizedToolQuery) &&
      !tool.description.toLowerCase().includes(normalizedToolQuery)
    ) {
      return [];
    }
    return [{ label, value: tool.id }];
  });

  async function handleToolsetsChange(items: MultiSelectPickerItem[]) {
    const nextToolsetIDs = items.map((item) => item.value);
    form.setValue("toolsets", nextToolsetIDs, {
      shouldDirty: true,
      shouldValidate: true,
    });

    const nextToolsetTools = collectToolsetTools(
      availableToolsets,
      nextToolsetIDs,
    );
    const nextDirectTools = form
      .getValues("tools")
      .filter((tool) => !nextToolsetTools.has(tool));
    form.setValue("tools", nextDirectTools, {
      shouldDirty: true,
      shouldValidate: true,
    });
  }

  async function handleToolsChange(items: MultiSelectPickerItem[]) {
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
          <Input
            {...form.register("name")}
            aria-label="Name"
            placeholder="Robot name"
          />
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Description</FormLabel>
          <Input
            {...form.register("description")}
            aria-label="Description"
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
            aria-label="Playbook"
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
            borderRadius="sm"
            resize="vertical"
            _focus={{
              outline: "none",
              borderColor: "accent.solid",
            }}
          />
          <FormErrorText>
            {form.formState.errors.playbook?.message}
          </FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Model</FormLabel>
          <ComboboxField
            control={form.control}
            name="model"
            items={models.map((model) => ({
              label: model.ref,
              value: model.ref,
            }))}
            placeholder={isCreating ? "Use default model" : "Select a model"}
            ariaLabel="Select model"
            disabled={!modelData || models.length === 0}
          />
          {modelError ? (
            <FormErrorText>Failed to load robot models.</FormErrorText>
          ) : (
            <Text variant="supporting">
              {isCreating
                ? "Leave unset to use the configured default model."
                : "Choose one of the enabled provider models."}
            </Text>
          )}
          <FormErrorText>{form.formState.errors.model?.message}</FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Toolsets</FormLabel>
          <MultiSelectPicker
            value={selectedToolsets}
            onChange={handleToolsetsChange}
            onQuery={setToolsetQuery}
            queryResults={filteredToolsetOptions}
            queryError={toolsetError ? "Failed to load Toolsets." : undefined}
            inputPlaceholder="Select Toolsets..."
            triggerProps={{ "aria-label": "Select Toolsets" }}
            size="sm"
          />
          <Text variant="supporting">
            Toolsets bundle reusable tools and specialist instructions. The
            Denbot can also discover Toolsets while it works.
          </Text>
        </FormControl>

        <FormControl>
          <FormLabel>Individual tools</FormLabel>
          <MultiSelectPicker
            value={selectedTools}
            onChange={handleToolsChange}
            onQuery={setToolQuery}
            queryResults={filteredToolOptions}
            queryError={toolError ? "Failed to load tools." : undefined}
            inputPlaceholder="Select tools..."
            triggerProps={{ "aria-label": "Select individual tools" }}
            size="sm"
          />
          <Text variant="supporting">
            Add one or two tools without assigning their entire Toolset. Tools
            already included by a selected Toolset are omitted automatically.
          </Text>
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
        intent="destructive"
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
            <Text variant="supporting" color="text.default">
              This will permanently delete this robot.
            </Text>
            <Text variant="supporting">
              Existing robot chat sessions will remain, but this robot will no
              longer be available.
            </Text>
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
              intent="destructive"
              variant="solid"
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
