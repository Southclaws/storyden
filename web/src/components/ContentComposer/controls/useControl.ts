import { useState } from "react";
import { useSlateStatic } from "slate-react";

import { Formats } from "../types";

import { isMarkActive, toggleMark } from "./utils";

export function useControl(format: Formats) {
  const editor = useSlateStatic();
  const [isActive, setIsActive] = useState(isMarkActive(editor, format));

  function onToggle() {
    const active = toggleMark(editor, format);
    setIsActive(active);
  }

  return {
    isActive,
    onToggle,
  };
}
