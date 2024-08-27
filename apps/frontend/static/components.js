const MAX_LENGTH_CHARS = 50

export function createTitleCard(title, link, className) {
  const titleCard = document.createElement("div")
  titleCard.setAttribute("class", className)
  titleCard.setAttribute("onclick", `window.open('${link}', 'mywindow')`)
  titleCard.textContent = title.length > MAX_LENGTH_CHARS ? title.slice(0, MAX_LENGTH_CHARS) + "..." : title

  return titleCard
}
