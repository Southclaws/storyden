import EmojiPicker, { EmojiStyle, type PickerProps } from "emoji-picker-react";

export default function ReactionPicker({
  onEmojiClick,
}: Pick<PickerProps, "onEmojiClick">) {
  return (
    <EmojiPicker
      onEmojiClick={onEmojiClick}
      emojiStyle={EmojiStyle.NATIVE}
      previewConfig={{ showPreview: false }}
      lazyLoadEmojis
    />
  );
}
