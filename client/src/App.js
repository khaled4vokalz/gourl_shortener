import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Shortener from "./components/shortener";

const App = () => (
  <Router>
    <Routes>
      <Route path="/" element={<Shortener />} />
    </Routes>
  </Router>
);

export default App;
