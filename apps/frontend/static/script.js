document.getElementById("search").addEventListener("submit", async function(event) {
  event.preventDefault()
  const resultsDiv = document.getElementById("search_results")
  const query = document.getElementById("query").value

  try {
    // TODO: Throw the URL of this endpoint in an environment variable
    const response = await fetch(
      `http://localhost:8000/videos?q=${encodeURIComponent(query)}`, 
      { 
        method: "GET",
        headers: {
          "Authorization": "Bearer 220fcc94-7a9a-45c2-bd1e-b54e70d6842e"
        }
      }
    )

    if (!response.ok) {
      console.log(`HTTP error, status ${response.status}`)
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
  resultItem.textContent = `${video.title} - ${video.channel_title}`

  return resultItem
}
