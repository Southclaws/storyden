import { Controller, ControllerProps } from "react-hook-form";

import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";

import { CancelAction } from "@/components/site/Action/Cancel";
import { SaveAction } from "@/components/site/Action/Save";
import { CardBox, HStack, WStack, styled } from "@/styled-system/jsx";

import { Byline } from "../../content/Byline";
import { ReactList } from "../ReactList/ReactList";
import { ReplyMenu } from "../ReplyMenu/ReplyMenu";

import { Form, Props, useReply } from "./useReply";

export function Reply(props: Props) {
  const { isEmpty, isEditing, resetKey, form, handlers } = useReply(props);

  const { thread, reply } = props;

  return (
    <CardBox id={reply.id}>
      <styled.form
        display="flex"
        flexDirection="column"
        gap="2"
        onSubmit={handlers.handleSave}
      >
        <WStack>
          <Byline
            href={`#${reply.id}`}
            author={reply.author}
            time={new Date(reply.createdAt)}
            updated={new Date(reply.updatedAt)}
          />

          {isEditing ? (
            <HStack>
              <>
                <CancelAction
                  type="button"
                  onClick={handlers.handleDiscardChanges}
                >
                  Discard
                </CancelAction>
                <SaveAction type="submit" disabled={isEmpty}>
                  Save
                </SaveAction>
              </>
            </HStack>
          ) : (
            <ReplyMenu
              thread={thread}
              reply={reply}
              onEdit={handlers.handleSetEditing}
            />
          )}
        </WStack>

        <ReplyBodyInput
          control={form.control}
          name="body"
          initialValue={reply.body}
          resetKey={resetKey}
          disabled={!isEditing}
          handleEmptyStateChange={handlers.handleEmptyStateChange}
        />
      </styled.form>

      <ReactList thread={thread} reply={reply} />
    </CardBox>
  );
}

type ReplyBodyInputProps = Omit<ControllerProps<Form>, "render"> & {
  initialValue: string;
  resetKey: string;
  handleEmptyStateChange: (isEmpty: boolean) => void;
};

function ReplyBodyInput({
  control,
  name,
  initialValue,
  resetKey,
  disabled,
  handleEmptyStateChange,
}: ReplyBodyInputProps) {
  return (
    <Controller<Form>
      render={({ field: { onChange } }) => {
        function handleChange(value: string, isEmpty: boolean) {
          handleEmptyStateChange(isEmpty);
          onChange(value);
        }

        return (
          <ContentComposer
            initialValue={initialValue}
            onChange={handleChange}
            resetKey={resetKey}
            disabled={disabled}
          />
        );
      }}
      control={control}
      name={name}
    />
  );
}
