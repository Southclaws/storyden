import { ChangeEvent, useEffect, useState } from "react";

export type Props = {
  initialValue: File | undefined;
  onSave: (f: File) => void;
};

export function useIconEditor(props: Props) {
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [file, setFile] = useState<File | undefined>(props.initialValue);

  useEffect(() => setFile(props.initialValue), [props.initialValue]);

  function onFileChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];

    if (!file) {
      throw new Error("Unexpected problem: File is missing from uploader.");
    }

    setFile(file);
  }

  function onSave() {
    if (file) props.onSave(file);
  }

  return { position, setPosition, onFileChange, onSave, file };
}
