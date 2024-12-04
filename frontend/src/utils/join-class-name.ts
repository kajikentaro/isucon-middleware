export function joinClassName(...classNames: (string | undefined | boolean)[]): string {
  return classNames.filter(Boolean).join(" ");
}