import { LogoutIcon } from "@/components/ui/icons/Logout";
import { LinkButtonStyleProps } from "@/components/ui/link-button";
import { Item } from "@/components/ui/menu";
import { API_ADDRESS } from "@/config";
import { button } from "@/styled-system/recipes";

import { AnchorProps } from "./Anchor";

// NOTE:
//
// Logging out in Storyden has a couple of steps, we need to allow logging out
// using standard HTTP without the need for JavaScript tricks. However, because
// Storyden can be run as a separate backend-frontend architecture where the API
// may run on a separate subdomain, the Clear-Site-Data header will not apply to
// requests made to a different origin. So if the Next.js app sets this header
// it will only apply to the Next.js origin, not the API origin. To solve this,
// logging out is handled by a form submission to the actual API itself, to
// the /auth/logout endpoint. This endpoint clears the site data such as cache
// and cookies for the API origin which resolves the caching for API calls such
// as the main endpoint used to check if the user is authenticated: /accounts
// But there's another problem, Next.js aggressively caches layout components
// including HTML and JS bundles for layouts and other components. This means we
// also need to tell Next.js to flush its caches. To achieve this, the API for
// logging out accepts a frontend path redirect parameter (not a full URL) which
// may be used to redirect post-logout. Here, we use it to redirect to a simple
// Next.js API route (/logout) which sets the Clear-Site-Data header for the
// frontend application origin (if it differs.) This ensures that both the API
// and frontend caches are cleared. Finally, the /logout route redirects back to
// the index page to complete the unfortunately complicated logout process.
//

export const LogoutID = "logout";
export const LogoutAction = `${API_ADDRESS}/api/auth/logout?redirect=${encodeURIComponent(`/logout`)}`;
export const LogoutLabel = "Logout";

type Props = AnchorProps & LinkButtonStyleProps;

export function LogoutAnchor({ hideLabel, ...props }: Props) {
  return (
    <form action={LogoutAction} method="POST">
      <button
        type="submit"
        className={button({ variant: "ghost", ...props })}
        title={LogoutLabel}
      >
        {<LogoutIcon />}
        {!hideLabel && (
          <>
            &nbsp;<span>{LogoutLabel}</span>
          </>
        )}
      </button>
    </form>
  );
}

export function LogoutMenuItem({ hideLabel }: AnchorProps) {
  return (
    <form action={LogoutAction} method="POST">
      <button type="submit" style={{ all: "unset", cursor: "pointer", width: "100%" }}>
        <Item value={LogoutID}>
          {<LogoutIcon />}
          {!hideLabel && (
            <>
              &nbsp;<span>{LogoutLabel}</span>
            </>
          )}
        </Item>
      </button>
    </form>
  );
}
