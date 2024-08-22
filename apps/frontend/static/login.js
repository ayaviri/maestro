import * as api from "./api.js"
import * as redirect from "./redirect.js"
import * as navbar from "./navbar.js"

document.getElementById("login").addEventListener("submit", async function(event) {
  event.preventDefault()
  const username = document.getElementById("username").value
  const password = document.getElementById("password").value

  await api.loginUser(username, password, redirect.toLoginPage)
})
