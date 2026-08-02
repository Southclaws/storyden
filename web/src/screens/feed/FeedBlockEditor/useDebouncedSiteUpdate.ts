import { debounce } from "lodash";
import { useEffect, useMemo, useRef } from "react";

import { AdminSettingsMutableProps } from "@/api/openapi-schema";

import { useFeedBlockEditor } from "./Context";

export function useDebouncedSiteUpdate() {
  const { updateSite } = useFeedBlockEditor();
  const updateSiteRef = useRef(updateSite);
  updateSiteRef.current = updateSite;

  const update = useMemo(
    () =>
      debounce(
        (patch: AdminSettingsMutableProps) => updateSiteRef.current(patch),
        500,
      ),
    [],
  );

  useEffect(() => {
    return () => {
      void update.flush();
    };
  }, [update]);

  return update;
}
