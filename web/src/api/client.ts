import { toast } from "sonner";

import { deriveError } from "@/utils/error";

import { Options, buildRequest, buildResult } from "./common";

export const fetcher = async <T>(opts: Options): Promise<T> => {
  const response = await fetch(buildRequest(opts));

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
