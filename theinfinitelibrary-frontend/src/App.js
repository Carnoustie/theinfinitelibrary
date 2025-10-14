import logo from './logotype.png';
import './App.css';
import { useState , useEffect, useRef} from 'react';
import { BrowserRouter, Routes, Route, Link , useNavigate, Navigate, useParams, useLocation} from 'react-router-dom';
import { redirect } from 'react-router';




function ChatRoom(props){

  const {chatId} = useParams()
  const [chatmessage, setChatMessage] = useState("")
  const [chatHistory, setChatHistory] = useState("")
  const esRef = useRef(null);


  useEffect(() => {
    console.log("\n\n\nhit useffect\n\n\n")
    const es = new EventSource(`http://localhost:8080/api/chatRoom/${chatId}`)
    esRef.current = es;
    es.onopen = () => console.log("SSE Open")
    es.onerror = (e) => console.log("SSE Error", e)


    es.onmessage = (ev) => {
      console.log(ev.data)
      console.log(chatHistory)
      setChatHistory(previousChatHistory => previousChatHistory + ev.data + "\n\n")
    }

    return ()=>{
      es.close()
    }
  }, [chatHistory])




  async function logToBackend(ev){
    ev.preventDefault()
    const response = await fetch(`/api/postMessage${chatId}`, {
      method: "POST",
      body: JSON.stringify({message: "Hello!"})
    })
  }

  async function submitHandler(ev){
    ev.preventDefault()
    console.log("\n\n\nhit here")
    const response = fetch(`/api/postMessage/${chatId}`, {
      method: "POST",
      body: chatmessage
    })
  }

  return(
    <div className="App-header">
      <form type="text" onSubmit={submitHandler}>
        <input 
        value = {chatmessage}
        onChange = {event => setChatMessage(event.target.value)}
        />
        <br/>
        <br/>
        <br/>
        <button type="submit" className="unplaced-button">
          post message
        </button>
      </form>
        <br/>
        <br/>
        <br/>
        <br/>
        <br/>
        <br/>
      <p className = "chathistory-style">{chatHistory}</p>

    </div>
  )
}




function SecurityInfo(props){
  
  const navigate = useNavigate();


  let returnButtonString = "Return to " + props.previousSite

  return(
    <div className="App-header">
      <h1 className="security-title-style">
        Security Information
      </h1>
      <p className="security-text-style">
        A text representation of your password is not stored in the database. Rather, a one-way encryption algorithm is applied to your password, the result of which is stored in the database. One-way means that decryption is intractible, making the encryption irreversible in order to inhibit recovery of the password text after encryption. Not storing the text representation in the database means that whoever can access the database, will not be able to access your password, and using the encrypted password can not be used to access your account. Additionally, a random salt is prepended to your password to further strengthen your password privacy (your encrypted password alone will thus be useless in attempts to access other sites using identical encryption). Overall, this approach is taken to preserve the safety of your account and the privacy of your password.
      </p>
      <button onClick={() => navigate(-1)} className="button-link">
        {returnButtonString}
      </button>
    </div>
  )
}



function AddBook(props){

  const [title, setTitle] = useState("")
  const [author, setAuthor] = useState("")
  const [addMessage, setAddMessage] = useState("")




  async function submitHandler(ev){
    ev.preventDefault()
    console.log("\n\n\nHit bookadder")
    const response = await fetch("api/addbook",{
      method: "POST",
      body: JSON.stringify({username: props.username, title: title, author:  author})
    });

    const r2 = await fetch("api/getbooks",{
      method: "POST",
      body: JSON.stringify({username:props.username})
    })
    
    const books = await r2.json()
    props.setBookList(books)
    console.log("\n\n\ntitle: ", title)
    
    var pm = title + " was added to your personal library!"
    console.log("\n\n\ntitle: ", title)
    
    setAddMessage(pm)
  }

  return(
    <div className="App-header">
      <Link to = "/Loggedin" className="button-link">return to profile</Link>
      <form className="form-style" onSubmit = {submitHandler}>
        <input
          value={title}
          onChange={event=>setTitle(event.target.value)}
          placeholder='Title'
        />
        <br/>
        <br/>
        <input
          value = {author}
          onChange={event=>setAuthor(event.target.value)}
          placeholder='Author'
        />
        <button type= "submit" className="button-style">
          Add Book
        </button>
      </form>
      <p>{addMessage}</p>
    </div>
  )
}


function Loggedin(props){
  console.log("Hit loggedin")
  return(
    <div className="App-header">
      <header>
        Welcome back {props.username}!
      </header>
      <Link to = "/addbook" className="button-style">
       Add book to your personal library
      </Link>       
      <p>
        You have read the following books:.... 
      </p>
      <ul className='book-list'>
        {props.bookList.length>0 ? props.bookList.map(b=>(
          <li key={b.title + b.author} className='book-item'> <p>{b.title} by {b.author}</p> 
            <Link to = {`/chatroom/${b.chatroom_id}`} className="lower-button-unfixed"> Enter book chat!</Link>
          </li>
        )) : null}
      </ul>
      <br/>
      <Link to = "/chatroom" className="lower-button-style">
        Enter chat room
      </Link>
    </div>
  )
}



function Signup(props){
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [message, setMessage] = useState("")


  props.setPreviousSite("signup page")


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
      <Link to= "/securityinfo" className="security-link-style">
          Security information.
      </Link>
      
    </div>
  );
}



function Login(props){
  // const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [message, setMessage] = useState("")
  const [isloggedin, setIsLoggedIn] = useState(false)

  props.setPreviousSite("Login page")

  async function loginHandler(ev){
    ev.preventDefault()

    const response = await fetch("/api/login", {
      method: "POST",
      body: JSON.stringify({username: props.username, password:  password})
    });

    setIsLoggedIn(true)
    const m = await response.text();
    setMessage(m)



    const r2 = await fetch("api/getbooks",{
      method: "POST",
      body: JSON.stringify({username:props.username})
    })

    const books = await r2.json()
    
    props.setBookList(books)



  }

  return(
    <div className="App-header">
      <header>
        <form className="form-style" onSubmit={loginHandler} >
          <input
            value={props.username}
            onChange={event=>props.setUsername(event.target.value)}
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
          {isloggedin ? <Navigate to = {`/loggedin`}/> : null}          
      </header>
      <Link to= "/securityinfo" className="security-link-style">
          Security information.
      </Link>
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

  const [bookList, setBookList] = useState([])
  const [previousSite, setPreviousSite] = useState("")
  const [chatRooms, setChatRooms] = useState([])


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
      element= {<Signup previousSite = {previousSite} setPreviousSite = {setPreviousSite}/>}
      />
      <Route 
      path="/login"
      element= {<Login username = {username} setUsername = {setUsername} bookList = {bookList} setBookList = {setBookList} previousSite = {previousSite} setPreviousSite = {setPreviousSite} chatRooms = {chatRooms} setChatRooms = {setChatRooms}/>}
      />
      <Route
      path="/loggedin"
      element=  {<Loggedin username = {username} setUsername= {setUsername} bookList = {bookList} setBookList = {setBookList}/>}
      />
      <Route
      path="/addbook"
      element= {<AddBook username = {username} setUsername = {setUsername} bookList = {bookList} setBookList = {setBookList}/>}
      />
      <Route
      path="/securityinfo"
      element = {<SecurityInfo previousSite = {previousSite} setPreviousSite = {setPreviousSite} />}
      />
      <Route
      path="/chatroom"
      element= {<ChatRoom/>}
      />
      <Route
      path="/chatroom/:chatId" element = {<ChatRoom/>}
      />
    </Routes>

    </BrowserRouter>
  );
}

export default App;