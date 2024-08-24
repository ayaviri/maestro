import * as utils from "./utils.js"

export function toLoginPage() {
  window.location.href = "/login.html"
  const errorElement = document.getElementsByClassName("error")[0]
  errorElement.textContent = `registration succeeded, but login failed :(`
}

export function unauthorisedUsers() {
  if (utils.getBearerToken() == null) {
    window.location.href = "/login.html"
  }
}
