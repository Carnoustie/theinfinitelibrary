import { useState, Dispatch, SetStateAction } from "react"

export type User = {
    username: string,
    password: string
}

export type Book = {
    title: string,
    author: string
}

export type ChatroomID = number
export type StringStateSetter = Dispatch<SetStateAction<string>>
export type NumberStateSetter = Dispatch<SetStateAction<number>>

export type LoginProps = {
    username: string;
    setUsername: StringStateSetter;
    bookList: Book[];
    setBookList: Dispatch<SetStateAction<Book[]>>;
    previousSite: string;
    setPreviousSite: StringStateSetter;
    chatrooms: ChatroomID[];
    setChatrooms: Dispatch<SetStateAction<ChatroomID[]>>;
}

export type UserProps = {
    username: string;
    setUsername: StringStateSetter;
    booklist: Book[];
    setBookList: Dispatch<SetStateAction<Book[]>>;
}

export type previousSiteProps = {
    previousSite: string;
    setPreviousSite: StringStateSetter;
}

export type ChatroomProps = {
    username: string;
    setUsername: StringStateSetter;
}