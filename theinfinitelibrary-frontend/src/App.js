import logo from './logotype.png';
import './App.css';
import { useState } from 'react';
import { BrowserRouter, Routes, Route, Link , useNavigate, Navigate, useParams} from 'react-router-dom';
import { redirect } from 'react-router';



function Loggedin(){
  console.log("Hit loggedin")
  const {username} = useParams()
  return(
    <div className="App-header">
      <header>
        Welcome back {username}!
      </header>
      <p>
        You have read the following books:....
      </p>
    </div>
  )
}



function Signup(){
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [message, setMessage] = useState("")


  async function submitHandler(ev){
    ev.preventDefault()

    const response = await fetch("/api/signup", {
      method: "POST",
      body: JSON.stringify({username, password})
    });

    const m = await response.text()
    setMessage(m)
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
      <pr className="vspace">
      Irreversible encryption is applied to your password to keep your account safe :) 
      </pr>
    </div>
  );
}



function Login(){
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [message, setMessage] = useState("")
  const [isloggedin, setIsLoggedIn] = useState(false)

  async function loginHandler(ev){
    ev.preventDefault()

    const response = await fetch("/api/login", {
      method: "POST",
      body: JSON.stringify({username,password})
    });
    setIsLoggedIn(true)
    const m = await response.text();
    setMessage(m)
    console.log("Hit here!")
  }

  return(
    <div className="App-header">
      <header>
        <form className="form-style" onSubmit={loginHandler} >
          <input
            value={username}
            onChange={event=>setUsername(event.target.value)}
            placeholder="Enter your username"
          />
          <br/>
          <br/>
          <input
            value = {password}
            onChange={event=>setPassword(event.target.value)}
            placeholder="Enter your password"
          />
          <br/>
          <button type = "submit" className="button-style">
            Login
          </button>
        </form>
          {message}
          {isloggedin ? <Navigate to = {`/loggedin/${username}`}/> : null}          
      </header>
      <pr className="vspace">
      Irreversible encryption is applied to your password to keep your account safe :) 
      </pr>
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
            <br/>
            <br/>
            <br/>
            <Link to="/login" className="button-link">
              Login!
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
  const [username, setUsername] = useState("")

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
      <Route 
      path="/login"
      element= <Login setUsername = {setUsername}/>
      />
      <Route
      path="/loggedin/:username"
      element= <Loggedin username = {username}/>
      />
    </Routes>

    </BrowserRouter>
  );
}

export default App;