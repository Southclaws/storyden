import { showUndoToast } from "@/components/ui/undo-toast";

type UndoableAction<T = void> = {
  action: () => Promise<T>;
  onUndo?: () => void;
  message: string;
  duration?: number;
  toastId?: string;
};

export async function withUndo<T = void>({
  action,
  onUndo,
  message,
  duration = 5000,
  toastId,
}: UndoableAction<T>): Promise<T | undefined> {
  return new Promise((resolve, reject) => {
    let isUndone = false;

    showUndoToast({
      message,
      duration,
      toastId,
      onUndo: () => {
        isUndone = true;
        onUndo?.();
        resolve(undefined);
      },
      onComplete: async () => {
        if (!isUndone) {
          try {
            const result = await action();
            resolve(result);
          } catch (error) {
            reject(error);
          }
        }
      },
    });
  });
}
