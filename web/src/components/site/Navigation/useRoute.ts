import { usePathname } from "next/navigation";

export type Route = {
  name: RouteName;
  label: string;
};

type RouteName = "index" | "library" | "admin" | "settings";

const routeLabels: Record<RouteName, string> = {
  index: "Home",
  library: "Library",
  admin: "Admin",
  settings: "Settings",
};

const mapping: Record<string, RouteName> = {
  "/": "index",
  "/l": "library",
  "/admin": "admin",
  "/settings": "settings",
};

function routeFromPrefix(prefix: string): Route | undefined {
  const routeName = mapping[prefix];
  if (!routeName) {
    return undefined;
  }
  return {
    name: routeName,
    label: routeLabels[routeName],
  };
}

export function useRoute(): Route | undefined {
  const pathname = usePathname();
  if (pathname[0] !== "/") {
    throw new Error(
      `useRoute: Invalid pathname "${pathname}". Expected a path starting with "/".`,
    );
  }

  const parts = pathname.split("/");
  if (parts.length < 2) {
    console.warn("useRoute: unexpected pathname format", pathname);
    return routeFromPrefix("/");
  }

  const first = parts[1];

  const prefix = `/${first}`;

  const route = routeFromPrefix(prefix);

  return route;
}
