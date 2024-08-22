import * as api from "./api.js"
import * as utils from "./utils.js"
import * as navbar from "./navbar.js"

if (utils.getBearerToken() == null) {
  window.location.href = "/login.html"
}

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

function createSearchResult(video) {
  const resultItem = document.createElement("div")
  resultItem.setAttribute("class", "search_result")
  resultItem.setAttribute("onclick", `window.open('${video.link}', 'mywindow')`)
  resultItem.textContent = `${decodeURIComponent(video.title)} - 
${decodeURIComponent(video.channel_title)}`
  console.log(video.title)

  return resultItem
}
