const ENDPOINT_URL = process.env.NEXT_PUBLIC_ENDPOINT_URL;

if (!ENDPOINT_URL) {
  throw new Error("env ENDPOINT_URL is not defined");
}

export function getReproduceUrl(ulid: string) {
  return `${ENDPOINT_URL}reproduce/${ulid}`;
}

export function getFetchAllUrl() {
  return `${ENDPOINT_URL}all`;
}

export type BodyType = "req-body" | "res-body" | "reproduced-res-body";

export function getBodyUrl(type: BodyType, ulid: string) {
  return `${ENDPOINT_URL}${type}/${ulid}`;
}
