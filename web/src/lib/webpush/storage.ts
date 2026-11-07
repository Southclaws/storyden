export type PushContext = {
  accountId: string;
  apiAddress: string;
  publicKey: string;
  subscriptionId?: string;
};

function openDatabase(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open("storyden-webpush", 1);

    request.onupgradeneeded = () => {
      request.result.createObjectStore("context");
    };

    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}

export async function readPushContext(): Promise<PushContext | undefined> {
  const database = await openDatabase();

  try {
    return await new Promise((resolve, reject) => {
      const transaction = database.transaction("context", "readonly");
      const request = transaction.objectStore("context").get("current");

      transaction.oncomplete = () => resolve(request.result);
      transaction.onabort = () => reject(transaction.error);
      transaction.onerror = () => reject(transaction.error);
    });
  } finally {
    database.close();
  }
}

export async function writePushContext(context?: PushContext): Promise<void> {
  const database = await openDatabase();

  try {
    await new Promise<void>((resolve, reject) => {
      const transaction = database.transaction("context", "readwrite");
      const store = transaction.objectStore("context");

      if (context) {
        store.put(context, "current");
      } else {
        store.delete("current");
      }

      transaction.oncomplete = () => resolve();
      transaction.onabort = () => reject(transaction.error);
      transaction.onerror = () => reject(transaction.error);
    });
  } finally {
    database.close();
  }
}
