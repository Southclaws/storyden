import { useEffect, useState } from "react";
import Markdown from "react-markdown";

import { Switch } from "@/components/ui/switch";
import { LStack, WStack, styled } from "@/styled-system/jsx";
import { htmlToMarkdown, markdownToHTML } from "@/utils/markdown";

type Props = {
  onChange: (value: string, isEmpty: boolean) => void;
  resetKey: string;
  initialValue?: string;
};

export function ContentComposerMarkdown(props: Props) {
  const [value, setValue] = useState(() => {
    if (props.initialValue) {
      return htmlToMarkdown(props.initialValue);
    }
    return "";
  });
  const [showPreview, setShowPreview] = useState(false);

  useEffect(() => {
    if (props.resetKey) {
      setValue("");
    }
  }, [props.resetKey]);

  async function onChange(e: React.ChangeEvent<HTMLTextAreaElement>) {
    const markdownRaw = e.target.value;

    setValue(markdownRaw);

    const html = await markdownToHTML(markdownRaw);

    const isEmpty = markdownRaw.trim().length === 0 || html.trim().length === 0;

    props.onChange(html, isEmpty);
  }

  function handleTogglePreview() {
    setShowPreview(!showPreview);
  }

  return (
    <LStack>
      <WStack justifyContent="end">
        <Switch size="sm" checked={showPreview} onClick={handleTogglePreview}>
          Preview
        </Switch>
      </WStack>
      {showPreview ? (
        <>
          <Markdown className="typography">
            {value || "_(Nothing to preview)_"}
          </Markdown>
        </>
      ) : (
        <styled.textarea
          onChange={onChange}
          value={value}
          lineHeight="relaxed"
          w="full"
          height="full"
          style={{
            height: "3lh",
          }}
        />
      )}
    </LStack>
  );
}
