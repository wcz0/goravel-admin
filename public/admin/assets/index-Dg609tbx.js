import{r as t,j as n,F as f,S as c}from"./index----IiH6k.js";function h({currentRoute:o}){const[a,s]=t.useState(!0),r=t.useRef(null),d=()=>{var i;const e=((i=document.querySelector(".ant-layout-content.overflow-hidden"))==null?void 0:i.offsetTop)||65;return`${window.innerHeight-e-75}px`},l=()=>{r.current.src=o.iframe_url};return t.useEffect(()=>{l();const e=()=>{r.current.style.height=d()};return window.addEventListener("resize",e),e(),()=>{window.removeEventListener("resize",e)}},[o.iframe_url]),n(f,{children:n(c,{spinning:a,style:{minHeight:a?"500px":""},children:n("iframe",{className:"owl-iframe",ref:r,title:"Iframe Page",width:"100%",style:{border:"none",order:0,boxSizing:"border-box",overflow:"hidden",minHeight:"500px"},onLoad:()=>{s(!1)}})})})}export{h as default};