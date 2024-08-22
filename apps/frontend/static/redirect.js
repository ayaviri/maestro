export function toLoginPage() {
  window.location.href = "/login.html"
  const errorElement = document.getElementsByClassName("error")[0]
  errorElement.textContent = `registration succeeded, but login failed :(`
}
