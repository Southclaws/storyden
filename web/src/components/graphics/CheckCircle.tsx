import { SVGProps } from "react";

export function CheckCircle(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      width="1em"
      height="1em"
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <g clip-path="url(#clip0_2773_13805)">
        <circle cx="12" cy="12" r="12" fill="#68D391" />
        <path
          fill-rule="evenodd"
          clip-rule="evenodd"
          d="M16.229 8.01796C16.6438 8.31427 16.7399 8.89078 16.4436 9.30562L11.8282 15.7672C11.6705 15.988 11.4236 16.1282 11.1532 16.1506C10.8828 16.1729 10.6162 16.0752 10.4244 15.8833L7.65513 13.1141C7.29464 12.7536 7.29464 12.1692 7.65513 11.8087C8.01561 11.4482 8.60007 11.4482 8.96056 11.8087L10.9593 13.8074L14.9413 8.23257C15.2376 7.81773 15.8141 7.72164 16.229 8.01796Z"
          fill="#303030"
        />
      </g>
      <defs>
        <clipPath id="clip0_2773_13805">
          <rect width="24" height="24" fill="white" />
        </clipPath>
      </defs>
    </svg>
  );
}
