import { ENV } from "@/constants";

const ENDPOINT_URL = ENV.ENDPOINT_URL;

export function getReproduceUrl(ulid: string) {
  return `${ENDPOINT_URL}isumid/reproduce/${ulid}`;
}

export function getSearchUrl(
  offset: number,
  length: number,
  query: string = ""
) {
  return `${ENDPOINT_URL}isumid/search?offset=${offset}&length=${length}&query=${query}`;
}

export type BodyType = "req-body" | "res-body" | "reproduced-res-body";

export function getBodyPath(type: BodyType, ulid: string) {
  if (process.env.NODE_ENV === "production") {
    // we don't need 'isumid' prefix as Next.js automatically add it according to basePath
    return `/${type}/${ulid}`;
  }

  return `${ENDPOINT_URL}isumid/${type}/${ulid}`;
}

export function getIsRecordingURL() {
  return `${ENDPOINT_URL}isumid/is-recording`;
}

export function getStartRecordingURL() {
  return `${ENDPOINT_URL}isumid/start-recording`;
}

export function getStopRecordingURL() {
  return `${ENDPOINT_URL}isumid/stop-recording`;
}

export function getRemoveAllURL() {
  return `${ENDPOINT_URL}isumid/remove-all`;
}

export function getRemoveURL(ulid: string) {
  return `${ENDPOINT_URL}isumid/remove/${ulid}`;
}
