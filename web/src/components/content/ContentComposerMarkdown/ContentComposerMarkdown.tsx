import { useEffect, useState } from "react";
import Markdown from "react-markdown";

import { Switch } from "@/components/ui/switch";
import { LStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  onChange: (value: string, isEmpty: boolean) => void;
  resetKey: string;
};

export function ContentComposerMarkdown(props: Props) {
  const [value, setValue] = useState("");
  const [showPreview, setShowPreview] = useState(false);

  // This is a huge hack but it means the composer doesn't need to be made into
  // a controlled component. Baiscally, if the resetKey changes, we reset the
  // content of the editor to the initial value or empty paragraph. Hacky? Yes.
  useEffect(() => {
    console.log("resetKey changed", props.resetKey);
    if (props.resetKey && value) {
      console.log("resetting to empty from", value);
      setValue("");
      return;
    }
  }, [setValue, props.resetKey]);

  async function onChange(e: React.ChangeEvent<HTMLTextAreaElement>) {
    const markdownRaw = e.target.value;

    setValue(markdownRaw);

    const html = /* convert markdownRaw to HTML */ markdownRaw;

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
