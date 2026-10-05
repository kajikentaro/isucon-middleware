// The built UI uses URLs relative to the page, so it works under whatever prefix the Go server is configured with.
// The dev server calls the Go server on :8080 directly (the APIs send Access-Control-Allow-Origin: *).
const BASE = import.meta.env.DEV ? "http://localhost:8080/isumid/" : "./";

export function getReproduceUrl(ulid: string) {
  return `${BASE}reproduce/${ulid}`;
}

export function getSearchUrl(
  offset: number,
  length: number,
  query: string = ""
) {
  return `${BASE}search?offset=${offset}&length=${length}&query=${query}`;
}

export type BodyType = "req-body" | "res-body" | "reproduced-res-body";

export function getBodyPath(type: BodyType, ulid: string) {
  return `${BASE}${type}/${ulid}`;
}

export function getIsRecordingURL() {
  return `${BASE}is-recording`;
}

export function getStartRecordingURL() {
  return `${BASE}start-recording`;
}

export function getStopRecordingURL() {
  return `${BASE}stop-recording`;
}

export function getRemoveAllURL() {
  return `${BASE}remove-all`;
}

export function getRemoveURL(ulid: string) {
  return `${BASE}remove/${ulid}`;
}

export function getExportURL(maxMetaMB: number, maxBodyMB: number) {
  return `${BASE}export?maxMetaMB=${maxMetaMB}&maxBodyMB=${maxBodyMB}`;
}

export function getExportSizeURL() {
  return `${BASE}export-size`;
}
