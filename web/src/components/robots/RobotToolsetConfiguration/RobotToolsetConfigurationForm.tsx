import { useState } from "react";

import { useRobotToolsList } from "@/api/openapi-client/robots";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormLabel } from "@/components/ui/form-label";
import { Input } from "@/components/ui/input";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/multi-select-picker";
import { Text } from "@/components/ui/text";
import { Textarea } from "@/components/ui/textarea";
import {
  RobotToolsetFormProps,
  useRobotToolsetConfigurationForm,
} from "@/screens/robots/useRobotToolsetConfigurationForm";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { useDisclosure } from "@/utils/useDisclosure";

export function RobotToolsetConfigurationForm(props: RobotToolsetFormProps) {
  const { form, isCreating, isEditable, handleSave } =
    useRobotToolsetConfigurationForm(props);
  const { data, error } = useRobotToolsList();
  const [query, setQuery] = useState("");
  const tools = data?.tools ?? [];
  const selectedToolIDs = form.watch("tools");
  const selectedTools: MultiSelectPickerItem[] = selectedToolIDs.map((id) => {
    const tool = tools.find((candidate) => candidate.id === id);
    return { label: tool?.name ?? id, value: id };
  });
  const normalizedQuery = query.trim().toLowerCase();
  const queryResults = tools.flatMap((tool) => {
    const label = tool.name ?? tool.id;

    if (
      !tool.available ||
      tool.toolset_only ||
      (normalizedQuery && !label.toLowerCase().includes(normalizedQuery))
    ) {
      return [];
    }

    return [{ label, value: tool.id }];
  });

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
        {!isEditable && props.toolset && (
          <HStack>
            <Badge size="sm" variant="outline">
              {props.toolset.source}
            </Badge>
            <Text variant="supporting">
              {props.toolset.source === "custom"
                ? "This Toolset is shared by another member. You can use it in your Robots, but only its author can edit it."
                : "This Toolset is supplied by Storyden or an installed plugin and is managed at its source."}
            </Text>
          </HStack>
        )}

        <FormControl>
          <FormLabel>Name</FormLabel>
          <Input
            {...form.register("name")}
            aria-label="Name"
            placeholder="Toolset name"
            disabled={!isEditable}
          />
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Description</FormLabel>
          <Input
            {...form.register("description")}
            aria-label="Description"
            placeholder="When this Toolset is useful"
            disabled={!isEditable}
          />
          <FormErrorText>
            {form.formState.errors.description?.message}
          </FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Instructions</FormLabel>
          <Textarea
            {...form.register("instruction")}
            aria-label="Instructions"
            placeholder="Specialist guidance applied whenever this Toolset is active..."
            rows={10}
            minH="40"
            disabled={!isEditable}
          />
          <Text variant="supporting">
            These instructions augment a Robot&apos;s playbook when the Toolset
            is assigned or loaded during a chat.
          </Text>
        </FormControl>

        <FormControl>
          <FormLabel>Tools</FormLabel>
          {isEditable ? (
            <MultiSelectPicker
              value={selectedTools}
              onChange={handleToolsChange}
              onQuery={setQuery}
              queryResults={queryResults}
              queryError={error ? "Failed to load tools." : undefined}
              inputPlaceholder="Select tools..."
              triggerProps={{ "aria-label": "Select Tools" }}
              size="sm"
            />
          ) : (
            <HStack gap="1" flexWrap="wrap">
              {selectedTools.length > 0 ? (
                selectedTools.map((tool) => (
                  <Badge key={tool.value} size="sm" variant="outline">
                    {tool.label}
                  </Badge>
                ))
              ) : (
                <Text variant="supporting">No directly exposed tools.</Text>
              )}
            </HStack>
          )}
        </FormControl>
      </LStack>

      {isEditable && (
        <WStack justifyContent="space-between">
          {props.onDelete ? (
            <ToolsetDeleteButton onDelete={props.onDelete} />
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
      )}
    </styled.form>
  );
}

function ToolsetDeleteButton({ onDelete }: { onDelete: () => Promise<void> }) {
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
        Delete Toolset
      </Button>

      <ModalDrawer
        title="Delete Toolset"
        isOpen={disclosure.isOpen}
        onClose={disclosure.onClose}
      >
        <LStack gap="6">
          <Text variant="supporting">
            This permanently deletes the Toolset. Toolsets assigned to a Robot
            must be removed from that Robot first.
          </Text>
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
              Delete Toolset
            </Button>
          </HStack>
        </LStack>
      </ModalDrawer>
    </>
  );
}
