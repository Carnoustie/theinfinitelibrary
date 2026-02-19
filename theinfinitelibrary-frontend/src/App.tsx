import logo from './resources/logotype.png';
import './style/App.css';
import { useState , useEffect, useRef} from 'react';
import { BrowserRouter, Routes, Route, Link , useNavigate, Navigate, useParams, useLocation} from 'react-router-dom';
import { redirect } from 'react-router';
import {Home, Login, Signup, Loggedin, AddBook, SecurityInfo, ChatRoom} from './components/Home'

//When deployed in container, backend URL from environment variable will be found

function App() {
  
  const [message, setMessage] = useState("");
  const [username, setUsername] = useState("")

  const [bookList, setBookList] = useState([])
  const [previousSite, setPreviousSite] = useState("")
  const [chatRooms, setChatRooms] = useState([])

  const join_prompt = "Join The Infinite Library!"
  return (
    <BrowserRouter>

    <Routes>
      <Route path="/"
      element={
        <Home
          join_prompt={join_prompt}
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
      path="/chatroom/:chatId" element = {<ChatRoom username={username} setUsername={setUsername}/>}
      />
    </Routes>

    </BrowserRouter>
  );
}

export default App;