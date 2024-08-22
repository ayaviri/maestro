import * as api from "./api.js"
import * as utils from "./utils.js"
import * as navbar from "./navbar.js"

if (utils.getBearerToken() == null) {
  window.location.href = "/login.html"
}

document.addEventListener("keyup", function(event) {
  if (event.key == "/") {
    document.getElementById("query").focus()
  }
})

document.getElementById("search").addEventListener("submit", async function(event) {
  event.preventDefault()
  const resultsDiv = document.getElementById("search_results")
  resultsDiv.textContent = ""
  const query = document.getElementById("query").value

  try {
    const response = await api.searchVideos(query)

    if (!response.ok) {
      const errorElement = document.getElementsByClassName("error")[0]
      errorElement.textContent = `:( (${response.status})`
    }

    const data = await response.json()

    if (data.videos.length > 0) {
      data.videos.forEach((video) => {
        resultsDiv.appendChild(createSearchResult(video))
      })
    } else {
      resultsDiv.textContent = "no results found :("
    }
  } catch (error) {
    console.log("There was a problem with the fetch operation: ", error)
  }
})

const MAX_LENGTH_CHARS = 50

function createSearchResult(video) {
  const container = document.createElement("div")
  container.setAttribute("class", "search_result")

  const resultTitle = document.createElement("div")
  resultTitle.setAttribute("class", "search_result_title")
  resultTitle.setAttribute("onclick", `window.open('${video.link}', 'mywindow')`)
  const title = `${decodeURIComponent(video.title)} - 
${decodeURIComponent(video.channel_title)}`
  resultTitle.textContent = title.length > MAX_LENGTH_CHARS ? title.slice(0, MAX_LENGTH_CHARS) + "..." : title

  const addToCart = document.createElement("div")
  addToCart.setAttribute("class", "add_to_cart")
  addToCart.setAttribute("id", video.id)
  addToCart.textContent = "+"
  addToCart.addEventListener("click", async function(event) {
    const videoId = event.target.id
    const response = await api.addToCart(videoId)

    if (!response.ok) {
      console.log(":(")
    } else {
      addToCart.textContent = ":)"
    }
  })

  container.appendChild(resultTitle)
  container.appendChild(addToCart)

  return container
}
