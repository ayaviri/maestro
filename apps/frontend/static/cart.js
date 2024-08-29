import * as navbar from "./navbar.js"
import * as api from "./api.js"
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

const cartItemFocusParams = {
  focusIndex: -1,
  items: document.getElementsByClassName("cart_item")
}
const menu = document.getElementById("menu")

document.addEventListener("keyup", function(event) {
  if (
    event.key == "j" && 
    !menu.classList.contains("active") && 
    cartItemFocusParams.items.length > 0
  ) {
    keybindings.focusNext(cartItemFocusParams)
  }

  if (
    event.key == "k" && 
    !menu.classList.contains("active") && 
    cartItemFocusParams.items.length > 0
  ) {
    keybindings.focusPrevious(cartItemFocusParams)
  }

  if (
    event.key == "o" && 
    document.activeElement &&
    document.activeElement.className == "cart_item"
  ) {
    const cartItemTitle = document.activeElement.getElementsByClassName("cart_item_title")[0]
    cartItemTitle.click()
  }

  if (
    event.key == "c" && 
    !event.shiftKey && 
    document.activeElement &&
    document.activeElement.className == "cart_item"
  ) {
    const removeFromCartButton = document.activeElement.getElementsByClassName("remove_from_cart")[0]
    removeFromCartButton.click()
  }

  if (event.key == "C") {
    const checkoutButton = document.getElementById("checkout")
    checkoutButton.click()
  }
})

document.addEventListener("DOMContentLoaded", async function(event) {
  event.preventDefault()
  const cartItemsDiv = document.getElementById("cart_items")
  const checkoutButton = document.getElementById("checkout")
  const errorElement = document.getElementById("error")
  let response

  try {
    response = await api.getCartItems()

  } catch (error) {
    errorElement.textContent = `could not connect to server :(`
    return
  }

  if (!response.ok) {
    errorElement.textContent = `could not load items from cart :(`
    return
  }

  const cartItems = (await response.json()).cart_items

  if (cartItems.length > 0) {
    cartItems.forEach((cartItem) => {
      cartItemsDiv.appendChild(createCartItem(cartItem))
    })
    checkoutButton.removeAttribute("hidden")
  } else {
    errorElement.textContent = "no items in cart :("
  }
})

document.getElementById("checkout").addEventListener("click", async function(event) {
  const checkoutButton = document.getElementById("checkout")
  checkoutButton.disabled = true
  const errorElement = document.getElementById("error")
  errorElement.textContent = ""
  let response

  try {
    response = await api.checkout()

    if (!response.ok) {
      errorElement.textContent = "could not request cart checkout:("
      return
    }
  } catch (error) {
    errorElement.textContent = "could not connect to server :("
    return
  }

  const jobId = (await response.json()).job_id
  const connection = api.openPersistentConnectionForDownloadStatus(jobId)
  const errorCallback = () => { 
    errorElement.textContent = "a song could not be downloaded :("
  }
  connection.addEventListener(
    "urls", async function(event) { 
      await api.downloadSongsUponCompletion(event, connection, errorCallback) 
      checkoutButton.disabled = false
    }
  )
})

function createCartItem(cartItem) {
  const container = document.createElement("div")
  container.setAttribute("class", "cart_item")
  container.setAttribute("tabindex", "0")
  const title = `${decodeURIComponent(cartItem.title)} - 
${decodeURIComponent(cartItem.channel_title)}`
  container.appendChild(
    components.createTitleCard(
      title,
      cartItem.link, 
      "cart_item_title"
    )
  )
  container.appendChild(createRemoveFromCartButton(cartItem))

  return container
}

function createRemoveFromCartButton(cartItem) {
  const removeFromCart = document.createElement("div")
  removeFromCart.setAttribute("class", "remove_from_cart")
  removeFromCart.setAttribute("id", cartItem.id)
  removeFromCart.textContent = "x"

  removeFromCart.addEventListener("click", async function(event) {
    const cartItemId = event.target.id
    const errorElement = document.getElementById("error")
    let response

    try {
      response = await api.removeFromCart(cartItemId)

      if (!response.ok) {
        errorElement.textContent = "could not remove item from cart :("
        return
      }
    } catch (error) {
      errorElement.textContent = "could not connect to server :("
      return
    }

    event.target.parentElement.remove()
    const cartItemCount = (await response.json()).remaining_count

    if (cartItemCount == 0) {
      const checkoutButton = document.getElementById("checkout")
      checkoutButton.setAttribute("hidden", "")
      errorElement.textContent = "no items in cart :("
    }
  })

  return removeFromCart
}
