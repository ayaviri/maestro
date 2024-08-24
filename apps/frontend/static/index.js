import * as api from "./api.js"
import * as utils from "./utils.js"
import * as navbar from "./navbar.js"
import * as redirect from "./redirect.js"


//   ___  _   _   ____   _    ____ _____   _     ___    _    ____  
//  / _ \| \ | | |  _ \ / \  / ___| ____| | |   / _ \  / \  |  _ \ 
// | | | |  \| | | |_) / _ \| |  _|  _|   | |  | | | |/ _ \ | | | |
// | |_| | |\  | |  __/ ___ \ |_| | |___  | |__| |_| / ___ \| |_| |
//  \___/|_| \_| |_| /_/   \_\____|_____| |_____\___/_/   \_\____/ 
//                                                                 

redirect.unauthorisedUsers()

//  _______     _______ _   _ _____ 
// | ____\ \   / / ____| \ | |_   _|
// |  _|  \ \ / /|  _| |  \| | | |  
// | |___  \ V / | |___| |\  | | |  
// |_____|  \_/  |_____|_| \_| |_|  
//                                  
//  _     ___ ____ _____ _____ _   _ _____ ____  ____  
// | |   |_ _/ ___|_   _| ____| \ | | ____|  _ \/ ___| 
// | |    | |\___ \ | | |  _| |  \| |  _| | |_) \___ \ 
// | |___ | | ___) || | | |___| |\  | |___|  _ < ___) |
// |_____|___|____/ |_| |_____|_| \_|_____|_| \_\____/ 
//                                                     

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

  const promises = [api.searchVideos(query), api.getCartItems()]
  let searchResponse, cartItemsResponse;
  [searchResponse, cartItemsResponse] = await Promise.all(promises)

  // if (!searchResponse.ok) {
  //   const errorElement = document.getElementsByClassName("error")[0]
  //   errorElement.textContent = `:( (${response.status})`
  //   return
  // }

  // if (!cartItemsResponse.ok) {
  // }

  const searchResults = (await searchResponse.json()).videos
  const cartItems = (await cartItemsResponse.json()).cart_items

  if (searchResults.length > 0) {
    searchResults.forEach((video) => {
      resultsDiv.appendChild(createSearchResult(video, cartItems))
    })
  } else {
    resultsDiv.textContent = "no results found :("
  }
})


//  _   _ _____ _     ____  _____ ____  ____  
// | | | | ____| |   |  _ \| ____|  _ \/ ___| 
// | |_| |  _| | |   | |_) |  _| | |_) \___ \ 
// |  _  | |___| |___|  __/| |___|  _ < ___) |
// |_| |_|_____|_____|_|   |_____|_| \_\____/ 
//                                            

const MAX_LENGTH_CHARS = 50
// These two constants represent the faces of the buttons
// for toggling a search result's presense in the user's cart
const IN_CART = ":)"
const NOT_IN_CART = "+"

function createSearchResult(video, cartItems) {
  const container = document.createElement("div")
  container.setAttribute("class", "search_result")
  container.appendChild(createSearchResultTitle(video))
  container.appendChild(createAddToCartButton(video, cartItems))

  return container
}

function createSearchResultTitle(video) {
  const resultTitle = document.createElement("div")
  resultTitle.setAttribute("class", "search_result_title")
  resultTitle.setAttribute("onclick", `window.open('${video.link}', 'mywindow')`)
  const title = `${decodeURIComponent(video.title)} - 
${decodeURIComponent(video.channel_title)}`
  resultTitle.textContent = title.length > MAX_LENGTH_CHARS ? title.slice(0, MAX_LENGTH_CHARS) + "..." : title

  return resultTitle
}

function createAddToCartButton(video, cartItems) {
  const addToCart = document.createElement("div")
  addToCart.setAttribute("class", "add_to_cart")
  addToCart.setAttribute("id", video.id)
  addToCart.textContent = isInCart(video, cartItems) ? IN_CART : NOT_IN_CART

  addToCart.addEventListener("click", async function(event) {
    const videoId = event.target.id
    const response = (
      addToCart.textContent == IN_CART ? 
      await api.removeFromCart(videoId) : 
      await api.addToCart(videoId)
    )

    if (!response.ok) {
      console.log(":(")
      return
    }


    addToCart.textContent = (
      addToCart.textContent == IN_CART ?
      NOT_IN_CART :
      IN_CART
    )
  })

  return addToCart
}

function isInCart(video, cartItems) {
  return cartItems.some((cartItem) => cartItem.id == video.id)
}
