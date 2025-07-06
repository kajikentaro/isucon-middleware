import { Header } from "../types";

export function makeCurlCommand(method: string, url: string, header: Header, body: string): string {
  const cmd: string[] = ["curl"];
  cmd.push("-X", method);
  // Headers
  for (const key in header) {
    for (const value of header[key]) {
      cmd.push("-H", `\"${key}: ${value}\"`);
    }
  }
  // Body
  if (body && body.length > 0) {
    cmd.push("--data-binary", `'${body}'`);
  }
 const fullUrl = new URL(url, window.location.origin);
  cmd.push(fullUrl.toString());
  return cmd.join(" ");
}
