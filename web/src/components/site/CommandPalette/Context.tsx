import {
  PropsWithChildren,
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";

export type CommandPaletteMode = "idle" | "chat";
export type CommandPaletteItem = "robot-chat" | "another";

type CommandPaletteContextValue = {
  open: boolean;
  setOpen: (open: boolean) => void;
  search: string;
  setSearch: (search: string) => void;
  handleSelectItem: (any) => Promise<void>;
  mode: CommandPaletteMode;
  setMode: (mode: CommandPaletteMode) => void;
  dialogRef: React.RefObject<HTMLDivElement | null>;
  focusInput: () => void;
};

const context = createContext<CommandPaletteContextValue | undefined>(
  undefined,
);

export function useCommandPalette() {
  const value = useContext(context);
  if (value === undefined) {
    throw new Error(
      "useCommandPalette must be used within a CommandPaletteProvider",
    );
  }
  return value;
}

export function CommandPaletteProvider({ children }: PropsWithChildren) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [mode, setMode] = useState<CommandPaletteMode>("idle");
  const dialogRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((open) => !open);
      }

      if (e.key === "backspace" && open && search === "") {
        e.preventDefault();
        setMode("idle");
      }

      if (e.key === "Escape" && open) {
        e.preventDefault();

        switch (mode) {
          case "chat":
            setMode("idle");
            break;

          case "idle":
            setOpen(false);
        }
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, [open, setOpen, mode, setMode, search]);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        dialogRef.current &&
        !dialogRef.current.contains(event.target as Node) &&
        // Only do outside click handling if the input is empty.
        search === ""
      ) {
        setOpen(false);
      }
    }

    if (open) {
      document.addEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [open, setOpen, search]);

  function focusInput() {
    const input = dialogRef.current?.querySelector(
      "[cmdk-input]",
    ) as HTMLInputElement | null;

    if (input) {
      input.focus();
      input.select();
    }
  }

  async function handleSelectItem(item: CommandPaletteItem) {
    switch (item) {
      case "robot-chat":
        setMode("chat");
        break;

      default:
        console.warn(`Unhandled command palette item: ${item}`);
        break;
    }
  }

  const value: CommandPaletteContextValue = {
    open,
    setOpen,
    search,
    setSearch,
    handleSelectItem,
    mode,
    setMode,
    dialogRef,
    focusInput,
  };

  return <context.Provider value={value}>{children}</context.Provider>;
}
