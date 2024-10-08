import * as utils from "./utils.js"

export async function registerUser(username, password) {
  const registrationForm = { username, password }
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
// 3) Redirects to the home page upon success
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
    throw new Error("login failed")
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

export async function addToCart(videoId) {
  const bearerToken = utils.getBearerToken()
  const response = await fetch(
    `http://localhost:8000/cart?v=${videoId}`, 
    { 
      method: "POST",
      headers: {
        "Authorization": `Bearer ${bearerToken}`
      }
    }
  )

  return response
}

export async function removeFromCart(videoId) {
  const bearerToken = utils.getBearerToken()
  const response = await fetch(
    `http://localhost:8000/cart?v=${videoId}`, 
    { 
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${bearerToken}`
      }
    }
  )

  return response
}

export async function getCartItems() {
  const bearerToken = utils.getBearerToken()
  const response = await fetch(
    `http://localhost:8000/cart`, 
    { 
      method: "GET",
      headers: {
        "Authorization": `Bearer ${bearerToken}`
      }
    }
  )

  return response
}

export async function checkout() {
  const bearerToken = utils.getBearerToken()
  const response = await fetch(
    `http://localhost:8000/checkout`,
    {
      method: "POST",
      headers: {
        "Authorization": `Bearer ${bearerToken}`,
      }
    }
  )

  return response
}

export function openPersistentConnectionForDownloadStatus(jobId) {
  return new EventSource(`http://localhost:8000/job/${jobId}`)
}

// Requests the given URL and downloads the response body as a blob
// using the filename in the Content-Disposition header
export async function downloadCart(url) {
  const bearerToken = utils.getBearerToken()
  const response = await fetch(url, {
    method: "GET",
    headers: {
      "Authorization": `Bearer ${bearerToken}`,
    }
  })

  console.log("after download request")

  if (!response.ok) {
    throw new Error("file not found")
  }

  const blobUrl = window.URL.createObjectURL(await response.blob())
  const downloadButton = document.createElement("a")
  downloadButton.id = blobUrl
  downloadButton.href = blobUrl
  downloadButton.download = _getFilenameFromResponseHeader(response)
  document.body.appendChild(downloadButton)
  downloadButton.click()
  downloadButton.remove()
}

function _getFilenameFromResponseHeader(response) {
  const contentDispositionHeader = response.headers.get("Content-Disposition")
  // "attachment; filename='file.ext'" => "file.ext"
  const filename = contentDispositionHeader.split(";")[1].trim().match("filename='(.*)'")[1]

  return filename
}
