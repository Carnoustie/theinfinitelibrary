import logo from './logotype.png';
import './App.css';
import { useState } from 'react';
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';




function Signup(){
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [message, setMessage] = useState("")


  function submitHandler(ev){
    ev.preventDefault()
    setMessage("Account handling and safe password handling yet to be implemented, your account was not saved.....")

    const response = fetch("/api/", {
      method: "POST",
      body: JSON.stringify({username, password})
    });

  }

  return(
    <div className="App-header">
      <header>
        <form className="form-style" onSubmit={submitHandler} >
          <input
            value={username}
            onChange={event=>setUsername(event.target.value)}
            placeholder="Choose a username"
          />
          <br/>
          <br/>
          <input
            value = {password}
            onChange={event=>setPassword(event.target.value)}
            placeholder="Choose a password"
          />
          <br/>
          <button type = "submit" className="button-style">
            Submit
          </button>
        </form>
          {message}
      </header>
    </div>
  );
}


function Home({message, join_prompt, logtoconsole}){
  return(
      <div>
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <br/>
          <nav>
            <Link to="/signup" className="button-link">
              Join The Infinite Library!
            </Link>
          </nav>
          <br/>
          {message && (
            <p className="response">
              {message}
            </p>
          )
          }
        </header>
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
      <Route 
      path="/signup"
      element= <Signup/>
      />
    </Routes>

    </BrowserRouter>
  );
}

export default App;