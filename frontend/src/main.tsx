import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App"; // We will create App.tsx next
import "./index.css"; // We will create index.css next

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
