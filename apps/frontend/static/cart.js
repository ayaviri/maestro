import * as navbar from "./navbar.js"
import * as api from "./api.js"

document.addEventListener("DOMContentLoaded", async function(event) {
  event.preventDefault()
  const cartItemsDiv = document.getElementById("cart_items")
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
  } else {
    cartItemsDiv.textContent = "no items in cart :("
  }
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
