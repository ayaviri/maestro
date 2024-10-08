import * as api from "./api.js"
import * as utils from "./utils.js"
import * as redirect from "./redirect.js"
import * as navbar from "./navbar.js"
import * as authForm from "./authForm.js"
import * as globalKeybindings from "./globalKeybindings.js"

document.getElementById("auth").addEventListener("submit", async function(event) {
  event.preventDefault()
  const username = document.getElementById("username").value
  const password = document.getElementById("password").value
  const errorElement = document.getElementById("error")
  let response

  try {
    response = await api.registerUser(username, password)

    if (!response.ok) {
      errorElement.textContent = "registration failed :("
      return
    }
  } catch (error) {
    errorElement.textContent = "could not connect to server :("
    return
  }


  try {
    await api.loginUser(username, password)
  } catch (error) {
    errorElement.textContent = "login failed :("
  }
})

