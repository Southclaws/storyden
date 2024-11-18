import { ChangeEvent, useEffect, useRef, useState } from "react";
import AvatarEditor from "react-avatar-editor";

export type Props = {
  initialValue: File | undefined;
  onSave: (f: Blob | null) => Promise<void>;
};

export function useIconEditor(props: Props) {
  const ref = useRef<AvatarEditor>(null);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [file, setFile] = useState<File | string>(props.initialValue ?? "");
  const [saving, setSaving] = useState(false);

  useEffect(
    () => props.initialValue && setFile(props.initialValue),
    [props.initialValue],
  );

  function onFileChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];

    if (!file) {
      throw new Error("Unexpected problem: File is missing from uploader.");
    }

    setFile(file);
  }

  function onSave() {
    if (!ref || !ref.current) {
      return;
    }

    setSaving(true);

    const canvasScaled = ref.current.getImageScaledToCanvas();
    canvasScaled.toBlob(async (f: Blob | null) => {
      await props.onSave(f);
      setSaving(false);
    });
  }

  return { ref, position, setPosition, onFileChange, onSave, saving, file };
}
