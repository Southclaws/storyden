import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useRef, useState } from "react";
import Markdown from "react-markdown";

import { IconButton } from "@/components/ui/icon-button";
import { EditIcon } from "@/components/ui/icons/Edit";
import { ShowIcon } from "@/components/ui/icons/ShowIcon";
import { Switch } from "@/components/ui/switch";
import { Box, HStack, LStack, styled } from "@/styled-system/jsx";
import { htmlToMarkdown, markdownToHTML } from "@/utils/markdown";

import { ContentComposerProps } from "../composer-props";

export function ContentComposerMarkdown(props: ContentComposerProps) {
  const [value, setValue] = useState(() => {
    if (props.initialValue) {
      return htmlToMarkdown(props.initialValue);
    }
    return "";
  });
  const [showPreview, setShowPreview] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    if (props.resetKey) {
      setValue("");
    }
  }, [props.resetKey]);

  useEffect(() => {
    if (props.disabled) return;

    const textarea = textareaRef.current;
    if (!textarea) return;

    const resizeObserver = new ResizeObserver(() => {
      textarea.style.height = "0px";
      textarea.style.height = `${textarea.scrollHeight}px`;
    });

    resizeObserver.observe(textarea.parentElement!);

    textarea.style.height = "0px";
    textarea.style.height = `${textarea.scrollHeight}px`;

    return () => {
      resizeObserver.disconnect();
    };
  }, [value, props.disabled]);

  async function onChange(e: React.ChangeEvent<HTMLTextAreaElement>) {
    const markdownRaw = e.target.value;

    setValue(markdownRaw);

    const html = await markdownToHTML(markdownRaw);

    const isEmpty = markdownRaw.trim().length === 0 || html.trim().length === 0;

    if (props.onChange) {
      props.onChange(html, isEmpty);
    }
  }

  function handleTogglePreview() {
    setShowPreview(!showPreview);
  }

  if (props.disabled) {
    return (
      <LStack position="relative" minHeight="8" flex="1">
        <Markdown className="typography">{value}</Markdown>
      </LStack>
    );
  }

  return (
    <LStack position="relative" minHeight="8" flex="1">
      <PreviewSwitch showPreview={showPreview} onChange={handleTogglePreview} />

      {showPreview ? (
        <>
          {value ? (
            <Markdown className="typography">{value}</Markdown>
          ) : (
            <styled.p height="14" color="fg.muted" fontStyle="italic">
              empty...
            </styled.p>
          )}
        </>
      ) : (
        <styled.textarea
          ref={textareaRef}
          onChange={onChange}
          value={value}
          lineHeight="relaxed"
          w="full"
          minHeight="0"
          resize="none"
          appearance="none"
          border="none"
          outline="none"
          color="fg.default"
          fontSize="md"
          transitionDuration="normal"
          transitionTimingFunction="default"
          _placeholder={{
            color: "fg.default",
          }}
          style={{
            border: "none",
            transitionProperty: "border-color, border-width",
            overflow: "hidden",
          }}
          placeholder="Write your heart out..."
        />
      )}
    </LStack>
  );
}

function PreviewSwitch({
  showPreview,
  onChange,
}: {
  showPreview: boolean;
  onChange: () => void;
}) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <Box
      position="absolute"
      right="0"
      p="1"
      opacity={isHovered ? "full" : "5"}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      cursor="pointer"
      backgroundColor="bg.subtle"
      backdropBlur="frosted"
      backdropFilter="auto"
      borderRadius="md"
      transition="all"
    >
      <HStack gap="2">
        <AnimatePresence>
          {isHovered && (
            <motion.div
              initial={{ width: 0, opacity: 0 }}
              animate={{ width: "auto", opacity: 1 }}
              exit={{ width: 0, opacity: 0 }}
              transition={{ duration: 0.2, ease: "easeInOut" }}
              style={{ overflow: "hidden" }}
            >
              <Switch size="sm" checked={showPreview} onClick={onChange}>
                Preview
              </Switch>
            </motion.div>
          )}
        </AnimatePresence>

        <IconButton type="button" variant="ghost" size="xs" onClick={onChange}>
          {showPreview ? <EditIcon w="4" /> : <ShowIcon w="4" />}
        </IconButton>
      </HStack>
    </Box>
  );
}
