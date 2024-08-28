const menu = document.getElementById("menu")
const menuFocusParams = {
  focusIndex: -1,
  items: document.getElementsByClassName("menu_link")
}

document.addEventListener("keyup", function(event) {

  if (event.key == "m" && noTextInputInFocus()) {
    menu.classList.toggle("active")

    if (menu.classList.contains("active")) {
      menuFocusParams.focusIndex = -1
      focusNext(menuFocusParams)
    }
  }

  if (event.key == "j" && menu.classList.contains("active")) {
    focusNext(menuFocusParams)
  }

  if (event.key == "k" && menu.classList.contains("active")) {
    focusPrevious(menuFocusParams)
  }
})

export function noTextInputInFocus() {
  return document.activeElement && 
    document.activeElement.tagName != "INPUT" && 
    document.activeElement.type != "text"
}

export function focusNext(params) {
  if (params.focusIndex < params.items.length - 1) {
    params.focusIndex++
  } else {
    params.focusIndex = 0
  }

  params.items[params.focusIndex].focus()
}

export function focusPrevious(params) {
  if (params.focusIndex == 0) {
    params.focusIndex = params.items.length - 1
  } else {
    params.focusIndex--
  }

  params.items[params.focusIndex].focus()
}
