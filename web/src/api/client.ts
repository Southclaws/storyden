import { toast } from "sonner";

import { deriveError } from "@/utils/error";

import { Options, buildRequest, buildResult } from "./common";

export const fetcher = async <T>(url: string, opts: Options): Promise<T> => {
  const request = buildRequest(url, {
    ...opts,
    // We use the browser default cache behaviour for the client side requests.
    // There's no revalidation set on the client as we're already using SWR for
    // that. The default cache behaviour will however make use of browser HTTP
    // Conditional Requests and ETag headers which some endpoints in Storyden
    // provide. This results in a mostly fast experience but it's slowed down a
    // bit by the server side behaviour (see server.ts comment for more info.)
    cache: "default",
  });

  const response = await fetch(request);

  return buildResult<T>(response);
};

type HandleArgs<T> = {
  onError?: (error: unknown) => Promise<void>;
  cleanup?: () => Promise<void>;
  promiseToast?: {
    loading: string;
    success: string;
  };
  action?: {
    label: string;
    onClick: (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => void;
  };
  errorToast?: boolean;
};

/**
 * Runs a UI action with shared error reporting and finally-style cleanup.
 *
 * Use `onError` for action-specific error handling and `cleanup` for work that
 * must run after either success or failure, such as clearing pending state.
 * These callbacks usually remove the need to wrap `handle` in another
 * try/catch/finally block.
 *
 * In the default mode, action errors are passed to `onError`, shown in an error
 * toast unless disabled, and handled by returning `undefined`. With
 * `promiseToast`, the toast manages loading/success/error feedback and the
 * returned promise rejects on failure; `onError` is not used in that mode.
 *
 * @param fn The asynchronous action to run.
 * @param args Optional error handling, toast configuration, and cleanup.
 * `cleanup` runs on both success and failure in either mode. Errors thrown by
 * `onError` or `cleanup` are not swallowed by the default mode.
 * @returns The action result on success, or `undefined` after a handled error
 * in the default mode. Promise-toast failures reject instead.
 */
export async function handle<T>(
  fn: () => Promise<T>,
  args?: HandleArgs<T>,
): Promise<T | undefined> {
  const {
    onError,
    cleanup,
    promiseToast,
    errorToast = true,
    action,
  } = args ?? {};

  if (promiseToast) {
    return new Promise((resolve, reject) => {
      toast.promise<T>(
        async () => {
          return await fn();
        },
        {
          loading: promiseToast.loading,
          success: (data: T) => {
            resolve(data);
            return promiseToast.success;
          },
          error: (error: unknown) => {
            reject(error);
            return deriveError(error);
          },
          finally: cleanup,
          action: action,
        },
      );
    });
  }

  try {
    return await fn();
  } catch (error) {
    await onError?.(error);

    if (errorToast) {
      toast.error(deriveError(error));
    }
  } finally {
    await cleanup?.();
  }

  return;
}
