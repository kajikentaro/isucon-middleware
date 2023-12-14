export function stringifyHeader(headerObj: {
  [key: string]: string[];
}): string {
  return Object.entries(headerObj)
    .map(([key, value]) => `${key}: ${value.join(", ")}`)
    .join("\n");
}
