import { uniq } from "lodash/fp";
import { type ChangeEvent, type KeyboardEvent, useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { Input } from "@/components/ui/input";
import { Flex, HStack, LStack } from "@/styled-system/jsx";

export type ModerationWordListEditorProps = {
  value: string[];
  onChange: (value: string[]) => void;
  disabled?: boolean;
};

export function ModerationWordListEditor({
  value,
  onChange,
  disabled,
}: ModerationWordListEditorProps) {
  const [newWord, setNewWord] = useState("");

  function handleNewWordChange(event: ChangeEvent<HTMLInputElement>) {
    setNewWord(event.target.value);
  }

  function handleNewWordSubmit() {
    const trimmed = newWord.trim();
    if (trimmed === "") return;
    onChange(uniq([...value, trimmed]));
    setNewWord("");
  }

  function handleNewWordKeyDown(event: KeyboardEvent<HTMLInputElement>) {
    if (event.key !== "Enter" || event.nativeEvent.isComposing) {
      return;
    }

    event.preventDefault();
    handleNewWordSubmit();
  }

  function handleRemoveWord(wordToRemove: string) {
    onChange(value.filter((word) => word !== wordToRemove));
  }

  return (
    <LStack>
      <Flex flexWrap="wrap">
        {value.map((word) => (
          <Badge key={word} pr="0">
            {word}
            <IconButton
              type="button"
              variant="ghost"
              aria-label={`Remove ${word}`}
              disabled={disabled}
              onClick={() => handleRemoveWord(word)}
            >
              <CancelIcon />
            </IconButton>
          </Badge>
        ))}
      </Flex>

      <HStack>
        <Input
          size="sm"
          value={newWord}
          disabled={disabled}
          onChange={handleNewWordChange}
          onKeyDown={handleNewWordKeyDown}
        />
        <Button
          type="button"
          size="sm"
          disabled={disabled}
          onClick={handleNewWordSubmit}
        >
          Add
        </Button>
      </HStack>
    </LStack>
  );
}
