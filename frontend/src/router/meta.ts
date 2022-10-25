export const getTitleTemplate = (
  title: string | undefined
) =>
  title
    ? `${title} - ${import.meta.env.VITE_APP_NAME}`
    : import.meta.env.VITE_APP_NAME;
