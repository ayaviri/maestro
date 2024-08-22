export const LS_BEARER_TOKEN = "bearer_token"
export function getBearerToken() {
  return localStorage.getItem(LS_BEARER_TOKEN)
}
