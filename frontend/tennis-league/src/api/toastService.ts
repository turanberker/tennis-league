type ToastHandler = (message: string) => void;

let toastFn: ToastHandler | null = null;

export const registerToast = (fn: ToastHandler) => {
  toastFn = fn;
};

export const showGlobalError = (message: string) => {
  toastFn?.(message);
};
