const ENDPOINT_URL = "http://localhost:8080"

export function getExecuteUrl(ulid: string){
    return `${ENDPOINT_URL}/execute/${ulid}`
}