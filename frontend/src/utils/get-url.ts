const ENDPOINT_URL = "http://localhost:8080";

export function getReproduceUrl(ulid: string) {
  return `${ENDPOINT_URL}/reproduce/${ulid}`;
}

export function getFetchAllUrl() {
  return `${ENDPOINT_URL}/all`;
}