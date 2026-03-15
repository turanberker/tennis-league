export const showGlobalError = (message: string) => {
  // Global bir event fırlatıyoruz
  const event = new CustomEvent('api-error', { detail: message });
  window.dispatchEvent(event);
};
