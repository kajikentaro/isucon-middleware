const ENDPOINT_URL = process.env.NEXT_PUBLIC_ENDPOINT_URL;

if (!ENDPOINT_URL) {
  throw new Error("env ENDPOINT_URL is not defined");
}

export function getReproduceUrl(ulid: string) {
  return `${ENDPOINT_URL}isumid/reproduce/${ulid}`;
}

export function getFetchListUrl(offset: number, length: number) {
  return `${ENDPOINT_URL}isumid/list?offset=${offset}&length=${length}`;
}

export type BodyType = "req-body" | "res-body" | "reproduced-res-body";

export function getBodyPath(type: BodyType, ulid: string) {
  if (process.env.NODE_ENV === "production") {
    // we don't need 'isumid' prefix as Next.js automatically add it according to basePath
    return `/${type}/${ulid}`;
  }

  return `${ENDPOINT_URL}isumid/${type}/${ulid}`;
}
