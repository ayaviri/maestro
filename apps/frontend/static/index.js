import * as api from "./api.js"
import * as utils from "./utils.js"
import * as navbar from "./navbar.js"
import * as redirect from "./redirect.js"
import * as components from "./components.js"
import * as keybindings from "./globalKeybindings.js"


//   ___  _   _   ____   ____ ____  ___ ____ _____   _     ___    _    ____  
//  / _ \| \ | | / ___| / ___|  _ \|_ _|  _ \_   _| | |   / _ \  / \  |  _ \ 
// | | | |  \| | \___ \| |   | |_) || || |_) || |   | |  | | | |/ _ \ | | | |
// | |_| | |\  |  ___) | |___|  _ < | ||  __/ | |   | |__| |_| / ___ \| |_| |
//  \___/|_| \_| |____/ \____|_| \_\___|_|    |_|   |_____\___/_/   \_\____/ 
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
    const queryInput = document.getElementById("query")

    if (keybindings.noTextInputInFocus()) {
      queryInput.focus()
    } else {
      const currentQuery = queryInput.value
      queryInput.value = currentQuery.substring(0, currentQuery.length - 1)
      queryInput.blur()
    }
  }
})

document.getElementById("search").addEventListener("submit", async function(event) {
  event.preventDefault()
  const resultsDiv = document.getElementById("search_results")
  resultsDiv.textContent = ""
  const errorElement = document.getElementById("error")
  errorElement.textContent = ""
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
    errorElement.textContent = "no results found :("
  }
})


//  _   _ _____ _     ____  _____ ____  ____  
// | | | | ____| |   |  _ \| ____|  _ \/ ___| 
// | |_| |  _| | |   | |_) |  _| | |_) \___ \ 
// |  _  | |___| |___|  __/| |___|  _ < ___) |
// |_| |_|_____|_____|_|   |_____|_| \_\____/ 
//                                            

// These two constants represent the faces of the buttons
// for toggling a search result's presense in the user's cart
const IN_CART = ":)"
const NOT_IN_CART = "+"

function createSearchResult(video, cartItems) {
  const container = document.createElement("div")
  container.setAttribute("class", "search_result")
  const title = `${decodeURIComponent(video.title)} - 
${decodeURIComponent(video.channel_title)}`
  container.appendChild(components.createTitleCard(title, video.link, "search_result_title"))
  container.appendChild(createAddToCartButton(video, cartItems))

  return container
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
