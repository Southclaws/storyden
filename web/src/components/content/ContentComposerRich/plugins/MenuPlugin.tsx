import { Editor } from "@tiptap/core";
import { EditorState, Plugin, PluginKey } from "@tiptap/pm/state";
import { EditorView } from "@tiptap/pm/view";
import { useCurrentEditor } from "@tiptap/react";
import { useEffect, useState } from "react";

import { css, cx } from "@/styled-system/css";

const PLUGIN_KEY = "floatingMenu";

const FloatingMenuPlugin = (options: FloatingMenuPluginProps) => {
  return new Plugin({
    key:
      typeof options.pluginKey === "string"
        ? new PluginKey(options.pluginKey)
        : options.pluginKey,
    view: (view) => new FloatingMenuView({ view, ...options }),
  });
};

type Optional<T, K extends keyof T> = Pick<Partial<T>, K> & Omit<T, K>;

type FloatingMenuProps = Omit<
  Optional<FloatingMenuPluginProps, "pluginKey" | "editor">,
  "element"
> & {
  children: React.ReactNode;
};

export interface FloatingMenuPluginProps {
  pluginKey: PluginKey | string;
  editor: Editor;
  element: HTMLElement;
  shouldShow?:
    | ((props: {
        editor: Editor;
        view: EditorView;
        state: EditorState;
        oldState?: EditorState;
      }) => boolean)
    | null;
}

export type FloatingMenuViewProps = FloatingMenuPluginProps & {
  view: EditorView;
};

export class FloatingMenuView {
  public editor: Editor;

  public element: HTMLElement;

  public view: EditorView;

  public preventHide = false;

  constructor({ editor, element, view }: FloatingMenuViewProps) {
    this.editor = editor;
    this.element = element;
    this.view = view;

    this.element.addEventListener("mousedown", this.mousedownHandler, {
      capture: true,
    });
    this.editor.on("focus", this.focusHandler);
    this.editor.on("blur", this.blurHandler);

    this.hide();
  }

  mousedownHandler = () => {
    this.preventHide = true;
  };

  focusHandler = () => {
    // we use `setTimeout` to make sure `selection` is already updated
    setTimeout(() => this.update(this.editor.view));
  };

  blurHandler = ({ event }: { event: FocusEvent }) => {
    if (this.preventHide) {
      this.preventHide = false;

      return;
    }

    if (
      event?.relatedTarget &&
      this.element.parentNode?.contains(event.relatedTarget as Node)
    ) {
      return;
    }

    this.hide();
  };

  update(view: EditorView, oldState?: EditorState) {
    if (!view.editable) {
      this.hide();
      return;
    }

    const { state } = view;
    const { doc, selection } = state;

    const isSame =
      oldState && oldState.doc.eq(doc) && oldState.selection.eq(selection);

    if (isSame) {
      return;
    }

    const { offsetX, offsetY } = this.getMenuOffset(view);

    this.element.setAttribute(
      "style",
      `top: ${offsetY}px; left: ${offsetX - 1}px;`,
    );
  }

  hide() {
    const style = this.element.getAttribute("style");
    this.element.setAttribute("style", `${style} visibility: hidden;`);
  }

  getMenuOffset(_view: EditorView) {
    // Position at top right of the editor container
    const offsetX = 0;
    const offsetY = 0;

    return { offsetX, offsetY };
  }
}

export const FloatingMenu = (props: FloatingMenuProps) => {
  const [element, setElement] = useState<HTMLDivElement | null>(null);
  const [isHovered, setIsHovered] = useState(false);
  const { editor: currentEditor } = useCurrentEditor();

  useEffect(() => {
    if (!element) {
      return;
    }

    if (props.editor?.isDestroyed || currentEditor?.isDestroyed) {
      return;
    }

    const menuEditor = props.editor || currentEditor;

    if (!menuEditor) {
      console.warn(
        "FloatingMenu component is not rendered inside of an editor component or does not have editor prop.",
      );
      return;
    }

    const plugin = FloatingMenuPlugin({
      pluginKey: PLUGIN_KEY,
      editor: menuEditor,
      element,
    });

    menuEditor.registerPlugin(plugin);
    return () => {
      menuEditor.unregisterPlugin(PLUGIN_KEY);
    };
  }, [props.editor, currentEditor, element]);

  return (
    <div
      ref={setElement}
      className={cx(menuStyles, isHovered ? hoverStyles : defaultStyles)}
      style={{ visibility: "hidden" }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {props.children}
    </div>
  );
};

const menuStyles = css({
  zIndex: "popover",
  position: "absolute",
  right: "0",
  top: "0",
  borderRadius: "md",
  display: "flex",
  flexWrap: "wrap",
  gap: "1",
  padding: "1",
  backgroundColor: "bg.subtle",
  backdropBlur: "frosted",
  backdropFilter: "auto",
  transition: "all",
  cursor: "pointer",
});

const defaultStyles = css({
  opacity: "5",
});

const hoverStyles = css({
  opacity: "full",
});
