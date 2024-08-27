import * as utils from "./utils.js"

class HamburgerMenu extends HTMLElement {
  connectedCallback() {
    this.innerHTML = `
 <div id="menu_toggle">
   <div id="hamburger"></div>
 </div>

 <nav id="menu">
   <ul>
     <li><a class="menu_link" href="/index.html">home</a></li>
     <li><a class="menu_link" href="/cart.html">cart</a></li>
     <li><a class="menu_link" href="/login.html">login</a></li>
     <li id="logout"><a class="menu_link" href="#">logout</a></li>
     <li><a class="menu_link" href="/register.html">register</a></li>
     <li><a class="menu_link" href="/keybindings.html">keybindings</a></li>
   </ul>
 </nav>
    `
  }
}

customElements.define("hamburger-menu", HamburgerMenu)
const menu = document.getElementById("menu")

document.getElementById("menu_toggle").addEventListener("click", function() {
  menu.classList.toggle("active")
})

document.getElementById("logout").addEventListener("click", utils.logout)
