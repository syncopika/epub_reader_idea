import './style.css';

import { 
  LoadEpubFile,
  NextFile,
  PrevFile,
} from '../wailsjs/go/main/App';

const loadFileButton = document.getElementById("load");

const iframe = document.getElementById("iframe");
iframe.contentWindow.document.body.innerHTML = "<h1 style='text-align: center; font-family: Arial'>this page is intentionally left blank 😃</h1>";

const prevPage = document.getElementById("prev");
const nextPage = document.getElementById("next");

loadFileButton.addEventListener('click', () => {
  LoadEpubFile();
});

prevPage.addEventListener('click', () => {
  PrevFile();
});

nextPage.addEventListener('click', () => {
  NextFile();
});

document.addEventListener('keydown', (evt) => {
  if(evt.code === 'ArrowLeft'){
    PrevFile();
  }else if(evt.code === 'ArrowRight'){
    NextFile();
  }
});

window.runtime.EventsOn('page', (pageData) => {
  // clear all links first to avoid issues when clicking on them
  // TODO: figure out a better way to deal with links?
  pageData = pageData.replaceAll(/href=".*"/g, "");
  
  // this is a bit hacky but since we're setting inner html of document.body,
  // we should remove elements like title and link elements since they don't belong in the body 
  // we'll add the styling separately
  const tempDiv = document.createElement('div');
  tempDiv.innerHTML = pageData;
  tempDiv.removeChild(tempDiv.querySelector('title'));
  tempDiv.querySelectorAll('link').forEach(link => {
    tempDiv.removeChild(link);
  });
  
  iframe.contentWindow.document.body.innerHTML = tempDiv.innerHTML;
});

window.runtime.EventsOn('style', (styleData) => {
  const styleElement = iframe.contentWindow.document.createElement('style');
  styleElement.textContent = styleData;
  iframe.contentWindow.document.head.appendChild(styleElement);
});