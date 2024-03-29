import { Button } from "src/theme/components/Button";

import { useControl } from "./useControl";

export function Italic() {
  const { isActive, onToggle } = useControl("italic");

  return (
    <Button
      backgroundColor={isActive ? "blackAlpha.100" : undefined}
      onMouseDown={onToggle}
      aria-label="Italic"
      type="button"
    >
      <svg
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M17.625 7.625V6.375H9.5V7.625H12.7125L9.98125 16.375H6.375V17.625H14.5V16.375H11.2875L14.0187 7.625H17.625Z"
          fill="#212529"
        />
      </svg>
    </Button>
  );
}
