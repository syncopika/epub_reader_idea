# epub_reader_idea

A super basic .epub reader made with Wails! It currently can display text, styling and some images (png and jpg).    
    
![gif of app](11-10-2024_193147.gif)    
You can find free eBooks at https://www.gutenberg.org/. This gif is showing the contents of *A New, Practical and Easy Method of Learning the Portuguese Language* by F. Ahn.    
    
You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
