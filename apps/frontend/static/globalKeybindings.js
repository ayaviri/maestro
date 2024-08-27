let currentFocusIndex = -1
const menu = document.getElementById("menu")
const menuLinks = document.getElementsByClassName("menu_link")

document.addEventListener("keyup", function(event) {
  if (event.key == "m" && noTextInputInFocus()) {
    menu.classList.toggle("active")

    if (menu.classList.contains("active")) {
      currentFocusIndex = -1
      focusNextItem()
    }
  }

  if (event.key == "j" && menu.classList.contains("active")) {
    focusNextItem()
  }

  if (event.key == "k" && menu.classList.contains("active")) {
    focusPreviousItem()
  }
})

export function noTextInputInFocus() {
  return document.activeElement && 
    document.activeElement.tagName != "INPUT" && 
    document.activeElement.type != "text"
}

function focusNextItem() {
  if (currentFocusIndex < menuLinks.length - 1) {
    currentFocusIndex++
  } else {
    currentFocusIndex = 0
  }

  menuLinks[currentFocusIndex].focus()
}

function focusPreviousItem() {
  if (currentFocusIndex == 0) {
    currentFocusIndex = menuLinks.length - 1
  } else {
    currentFocusIndex--
  }

  menuLinks[currentFocusIndex].focus()
}
