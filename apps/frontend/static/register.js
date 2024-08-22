import * as api from "./api.js"
import * as utils from "./utils.js"
import * as redirect from "./redirect.js"
import * as navbar from "./navbar.js"

document.getElementById("register").addEventListener("submit", async function(event) {
  event.preventDefault()
  const email = document.getElementById("email").value
  const username = document.getElementById("username").value
  const password = document.getElementById("password").value

  let response = await api.registerUser(email, username, password)

  if (!response.ok) {
    const errorMessage = await response.text()
    errorElement.textContent = `registration failed :(`
    return
  }

  await api.loginUser(username, password, redirect.toLoginPage)
})

