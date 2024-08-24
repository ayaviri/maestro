import * as utils from "./utils.js"

class HamburgerMenu extends HTMLElement {
  connectedCallback() {
    this.innerHTML = `
 <div id="menu_toggle">
   <div id="hamburger"></div>
 </div>

 <nav id="menu">
   <ul>
     <li><a href="/index.html">home</a></li>
     <li><a href="/cart.html">cart</a></li>
     <li><a href="/login.html">login</a></li>
     <li id="logout"><a href="#">logout</a></li>
     <li><a href="/register.html">register</a></li>
   </ul>
 </nav>
    `
  }
}

customElements.define("hamburger-menu", HamburgerMenu)
const menu = document.getElementById("menu")
const hamburger = document.getElementById("hamburger")

document.getElementById("menu_toggle").addEventListener("click", function() {
  menu.classList.toggle("active")
  hamburger.classList.toggle("active")
})

document.getElementById("logout").addEventListener("click", utils.logout)
