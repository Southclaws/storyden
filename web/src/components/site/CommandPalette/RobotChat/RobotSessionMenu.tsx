"use client";

import { useCombobox, useFilter, useListCollection } from "@ark-ui/react";
import { useEffect, useMemo } from "react";

import { RobotSessionRef } from "@/api/openapi-schema";
import * as Combobox from "@/components/ui/combobox";
import { IconButton } from "@/components/ui/icon-button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { ChevronUpDownIcon } from "@/components/ui/icons/Chevron";
import { Input } from "@/components/ui/input";

import { useCommandPalette } from "../Context";

import { useRobotChat } from "./RobotChatContext";

type Item = {
  label: string;
  value: string;
};

function sessionToItem(session: RobotSessionRef): Item {
  return {
    label: session.name,
    value: session.id,
  };
}

export function RobotSessionMenu() {
  const { sessions, sessionId } = useRobotChat();
  const { loadChatSession } = useCommandPalette();

  const { contains } = useFilter({ sensitivity: "base" });

  const items = useMemo(() => sessions.map(sessionToItem), [sessions]);

  const { collection, filter } = useListCollection({
    initialItems: items,
    itemToString: (item) => item.label,
    itemToValue: (item) => item.value,
    filter: contains,
  });

  const handleInputChange = (details: Combobox.InputValueChangeDetails) => {
    filter(details.inputValue);
  };

  function handleChange({ value }: Combobox.ValueChangeDetails) {
    if (!value || value.length === 0) {
      return;
    }

    const selectedSessionId = value[0];
    if (selectedSessionId && selectedSessionId !== sessionId) {
      loadChatSession(selectedSessionId);
    }
  }

  const combobox = useCombobox({
    collection,
    value: [sessionId],
    onInputValueChange: handleInputChange,
    onValueChange: handleChange,
  });

  const currentSessionName =
    sessions.find((s) => s.id === sessionId)?.name || "Current Session";

  return (
    <Combobox.RootProvider value={combobox} size="xs">
      <Combobox.Control>
        <Combobox.Input placeholder={currentSessionName} asChild>
          <Input size="xs" />
        </Combobox.Input>
        <Combobox.Trigger asChild>
          <IconButton variant="link" aria-label="Select session" size="xs">
            <ChevronUpDownIcon />
          </IconButton>
        </Combobox.Trigger>
      </Combobox.Control>
      <Combobox.Positioner>
        <Combobox.Content>
          <Combobox.List>
            <Combobox.ItemGroup>
              <Combobox.ItemGroupLabel>Recent Sessions</Combobox.ItemGroupLabel>
              {collection.items.map((item) => (
                <Combobox.Item key={item.value} item={item}>
                  <Combobox.ItemText>{item.label}</Combobox.ItemText>
                  <Combobox.ItemIndicator>
                    <CheckIcon />
                  </Combobox.ItemIndicator>
                </Combobox.Item>
              ))}
            </Combobox.ItemGroup>
          </Combobox.List>
        </Combobox.Content>
      </Combobox.Positioner>
    </Combobox.RootProvider>
  );
}
