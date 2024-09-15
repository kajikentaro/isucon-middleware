export const MAX_ROW_LENGTH = 100;

const envTmp = {
  TOP_PAGE_PATH: process.env.NEXT_PUBLIC_TOP_PAGE_PATH,
  ENDPOINT_URL: process.env.NEXT_PUBLIC_ENDPOINT_URL,
};
for (const [key, value] of Object.entries(envTmp)) {
  if (typeof value === "undefined") {
    throw new Error(`env ${key} is not defined`);
  }
}

export const ENV = envTmp as Record<string, string>;
