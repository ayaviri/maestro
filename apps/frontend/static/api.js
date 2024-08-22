import * as utils from "./utils.js"

export async function registerUser(email, username, password) {
  const registrationForm = { email, username, password }
  const response = await fetch(
    `http://localhost:8000/register`,
    {
      method: "POST",
      body: JSON.stringify(registrationForm),
    }
  )

  return response
}

// 1) Hits the /login endpoint on the backend
// 2) Places the bearer token in the response body into localStorage
// 3a) Calls the given error callback upon failure
// 3b) Or redirects to the home page upon success
export async function loginUser(username, password, errorCallback) {
  const loginForm = { username, password }
  const response = await fetch(
    `http://localhost:8000/login`,
    {
      method: "POST",
      body: JSON.stringify(loginForm),
    }
  )

  if (!response.ok) {
    errorCallback()
  } else {
    const data = await response.json()
    localStorage.setItem(utils.LS_BEARER_TOKEN, data.bearer_token)
    window.location.href = "/index.html"
  }
}

export async function searchVideos(query) {
  const bearerToken = utils.getBearerToken()
  const response = await fetch(
    `http://localhost:8000/videos?q=${encodeURIComponent(query)}`, 
    { 
      method: "GET",
      headers: {
        "Authorization": `Bearer ${bearerToken}`
      }
    }
  )

  return response
}
