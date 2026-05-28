package com.example.TIL_Springboot_Backend;

import jakarta.persistence.EntityNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class BookService {

    private BookRepository bookRepository;
    private UserRepository userRepository;

    public BookService(BookRepository bookRepository, UserRepository userRepository){
        this.bookRepository = bookRepository;
        this.userRepository = userRepository;
    }

    public Book addBookToUser(String title, String author, Long userId){
        Book book = new Book(title, author, userId);
        userRepository.findById(userId);
        if(!userRepository.existsById(userId)){
            throw new EntityNotFoundException(userId.toString());
        }
        return bookRepository.save(book);
    }
}
