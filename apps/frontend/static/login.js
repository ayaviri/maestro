import * as api from "./api.js"
import * as redirect from "./redirect.js"
import * as navbar from "./navbar.js"
import * as authForm from "./authForm.js"
import * as globalKeybindings from "./globalKeybindings.js"

document.getElementById("auth").addEventListener("submit", async function(event) {
  event.preventDefault()
  const username = document.getElementById("username").value
  const password = document.getElementById("password").value

  try {
    await api.loginUser(username, password)
  } catch (error) {
    const errorElement = document.getElementById("error")
    errorElement.textContent = "login failed :("
  }
})
