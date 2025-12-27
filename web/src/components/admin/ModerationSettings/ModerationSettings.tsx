import { uniq } from "lodash/fp";
import { ChangeEvent, useState } from "react";
import { Controller } from "react-hook-form";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { NumberInputField } from "@/components/ui/form/NumberInputField";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { Input } from "@/components/ui/input";
import {
  CardBox,
  Flex,
  HStack,
  LStack,
  WStack,
  styled,
} from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Props, useModerationSettings } from "./useModerationSettings";

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
      <CardBox className={lstack()}>
        <WStack>
          <Heading size="md">Moderation settings</Heading>
          <Button type="submit" loading={formState.isSubmitting}>
            Save
          </Button>
        </WStack>

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

            <Controller
              control={control}
              name="wordReportList"
              render={({ field, fieldState, formState }) => {
                const [newWord, setNewWord] = useState("");

                function handleNewWordChange(v: ChangeEvent<HTMLInputElement>) {
                  setNewWord(v.target.value);
                }

                function handleNewWordSubmit() {
                  const trimmed = newWord.trim();
                  if (trimmed === "") return;
                  const newList = uniq([...(field.value ?? []), trimmed]);
                  field.onChange(newList);
                  setNewWord("");
                }

                function handleRemoveWord(wordToRemove: string) {
                  const newList =
                    field.value?.filter((w) => w !== wordToRemove) ?? [];
                  field.onChange(newList);
                }

                return (
                  <LStack>
                    <Flex flexWrap="wrap">
                      {field.value?.map((word) => (
                        <Badge key={word} pr="0">
                          {word}
                          <IconButton
                            type="button"
                            size="xs"
                            variant="ghost"
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
                        onChange={handleNewWordChange}
                      />
                      <Button
                        type="button"
                        size="sm"
                        onClick={handleNewWordSubmit}
                      >
                        Add
                      </Button>
                    </HStack>
                  </LStack>
                );
              }}
            />

            <FormHelperText>
              Words and phrases that will automatically report and hide posts
              that contain them.
            </FormHelperText>
          </FormControl>
          <FormControl>
            <FormLabel>Word block list</FormLabel>

            <Controller
              control={control}
              name="wordBlockList"
              render={({ field, fieldState, formState }) => {
                const [newWord, setNewWord] = useState("");

                function handleNewWordChange(v: ChangeEvent<HTMLInputElement>) {
                  setNewWord(v.target.value);
                }

                function handleNewWordSubmit() {
                  const trimmed = newWord.trim();
                  if (trimmed === "") return;
                  const newList = uniq([...(field.value ?? []), trimmed]);
                  field.onChange(newList);
                  setNewWord("");
                }

                function handleRemoveWord(wordToRemove: string) {
                  const newList =
                    field.value?.filter((w) => w !== wordToRemove) ?? [];
                  field.onChange(newList);
                }

                return (
                  <LStack>
                    <Flex flexWrap="wrap">
                      {field.value?.map((word) => (
                        <Badge key={word} pr="0">
                          {word}
                          <IconButton
                            type="button"
                            size="xs"
                            variant="ghost"
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
                        onChange={handleNewWordChange}
                      />
                      <Button
                        type="button"
                        size="sm"
                        onClick={handleNewWordSubmit}
                      >
                        Add
                      </Button>
                    </HStack>
                  </LStack>
                );
              }}
            />

            <FormHelperText>
              Words and phrases that will instantly reject posts without
              creating a report.
            </FormHelperText>
          </FormControl>
        </Flex>
      </CardBox>
    </styled.form>
  );
}
