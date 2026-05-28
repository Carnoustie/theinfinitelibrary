package com.example.TIL_Springboot_Backend;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;

@Entity
public class Book {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long bookId;

    private String title;
    private String author;
    private Long userId;

    public Book(){ }

    public Book(String title, String author, Long userId){
        this.title = title;
        this.author = author;
        this.userId = userId;
    }

    public Long getBookId() {
        return bookId;
    }

    public Long getUserId(){ return userId; }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getTitle() {
        return title;
    }

    public void setAuthor(String author) {
        this.author = author;
    }

    public String getAuthor() {
        return author;
    }
}