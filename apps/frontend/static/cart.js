import * as navbar from "./navbar.js"
import * as api from "./api.js"
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

document.addEventListener("DOMContentLoaded", async function(event) {
  event.preventDefault()
  const cartItemsDiv = document.getElementById("cart_items")
  const checkoutButton = document.getElementById("checkout")
  const response = await api.getCartItems()

  if (!response.ok) {
    console.log(":(")
    return
  }

  const cartItems = (await response.json()).cart_items

  if (cartItems.length > 0) {
    cartItems.forEach((cartItem) => {
      cartItemsDiv.appendChild(createCartItem(cartItem))
    })
    checkoutButton.removeAttribute("hidden")
  } else {
    cartItemsDiv.textContent = "no items in cart :("
  }
})

document.getElementById("checkout").addEventListener("click", async function(event) {
  const response = await api.checkout(() => console.log("checkout failed"))
})

function createCartItem(cartItem) {
  const container = document.createElement("div")
  container.setAttribute("class", "cart_item")
  container.appendChild(createCartItemTitle(cartItem))
  container.appendChild(createRemoveFromCartButton(cartItem))

  return container
}

const MAX_LENGTH_CHARS = 50

function createCartItemTitle(cartItem) {
  const cartItemTitle = document.createElement("div")
  cartItemTitle.setAttribute("class", "cart_item_title")
  cartItemTitle.setAttribute("onclick", `window.open('${cartItem.link}', 'mywindow')`)
  const title = `${decodeURIComponent(cartItem.title)} - 
${decodeURIComponent(cartItem.channel_title)}`
  cartItemTitle.textContent = title.length > MAX_LENGTH_CHARS ? title.slice(0, MAX_LENGTH_CHARS) + "..." : title

  return cartItemTitle
}

function createRemoveFromCartButton(cartItem) {
  const removeFromCart = document.createElement("div")
  removeFromCart.setAttribute("class", "remove_from_cart")
  removeFromCart.setAttribute("id", cartItem.id)
  removeFromCart.textContent = "x"

  removeFromCart.addEventListener("click", async function(event) {
    const cartItemId = event.target.id
    const response = await api.removeFromCart(cartItemId)

    if (!response.ok) {
      console.log(":(")
      return
    }

    event.target.parentElement.remove()
  })

  return removeFromCart
}
