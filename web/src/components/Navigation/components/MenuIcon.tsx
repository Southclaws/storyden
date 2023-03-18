import { SVGProps } from "react";

export function MenuIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      width={64}
      height={22}
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        d="M12.666 5.445h38.667M12.666 16.556h38.667M12.666 11h38.667"
        stroke="#000"
        strokeOpacity={0.16}
        strokeWidth={1.5}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}
