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
     <li><a href="/login.html">login</a></li>
     <li id="logout"><a href="#">logout</a></li>
     <li><a href="/register.html">register</a></li>
     <li><a href="#">cart</a></li>
   </ul>
 </nav>
    `
  }
}

customElements.define("hamburger-menu", HamburgerMenu)

document.getElementById("menu_toggle").addEventListener("click", function() {
  document.getElementById("menu").classList.toggle("active")
  document.getElementById("hamburger").classList.toggle("active")
})

document.getElementById("logout").addEventListener("click", utils.logout)
