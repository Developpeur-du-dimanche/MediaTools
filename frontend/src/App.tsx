import { HashRouter, Routes, Route } from "react-router-dom";
import HomePage from './pages/Home';

// Show the actual Wails version this project was generated against.
function App() {


  return (
    <HashRouter basename={"/"}>
      {/* The rest of your app goes here */}
      <Routes>
        <Route path="/" element={<HomePage />} />
        {/* more... */}
      </Routes>
    </HashRouter>
  );
}

export default App
