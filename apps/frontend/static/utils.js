export const LS_BEARER_TOKEN = "bearer_token"
export function getBearerToken() {
  return localStorage.getItem(LS_BEARER_TOKEN)
}

export function logout() {
  localStorage.removeItem(LS_BEARER_TOKEN)
  window.location.href = "/index.html"
}
