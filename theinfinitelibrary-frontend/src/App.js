import logo from './logotype.png';
import './App.css';
import { useState } from 'react';
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';



function About(){
  return( 
    <div className="App-header">
      <h1>The infinite Library is..........</h1>
    </div>
  )
}


function Home({message, join_prompt, logtoconsole}){
  return(
      <div>
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <button className="join_button" type="button" onClick={logtoconsole}>
            {join_prompt}
          </button>
          <br/>
          <br/>
          {message && (
            <p className="response">
              {message}
            </p>
          )
          }
        </header>
        <footer className="footer-style">
          <nav>
            <Link to="/about">About</Link>
          </nav>
        </footer>
      </div>
  );
}




function App() {
  
  const [message, setMessage] = useState("");

  function logtoconsole(){
    console.log("\n\n\nSomebody wants to join the book club!\n\n\n")
    fetch("/api/").then(res => res.text()).then(data => {
      console.log(data);
      setMessage(data);
    })
  }

  const join_prompt = "Join The Infinite Library!"
  return (
    <BrowserRouter>

    <Routes>
      <Route path="/"
      element={
        <Home
          message={message}
          join_prompt={join_prompt}
          logtoconsole={logtoconsole}
        />
        }
      />
      <Route path="/about" element= {<About/>} />
    </Routes>

    </BrowserRouter>
  );
}

export default App;