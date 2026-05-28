package com.example.TIL_Springboot_Backend;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;
import com.example.TIL_Springboot_Backend.DTOs.*;

@RestController
public class BookController {

    @Autowired
    private BookRepository bookRepository;
    @Autowired
    private UserRepository userRepository;


    @PostMapping("api/addBook")
    public ResponseEntity<Book> addBook(@RequestBody BookRequest b){
        //StrictHttpFirewall fw = new StrictHttpFirewall();
        System.out.println("\n\n\nhit here");

        Book bk = new Book(b.title(), b.author(), b.userId());
        bookRepository.save(bk);
        System.out.println(bk.getTitle());
        return ResponseEntity.ok(bk);
    }
    /*
    @GetMapping("api/getbooks")
    public ResponseEntity<List<Book>> getBooks(@RequestBody UserRequest request){
        List<User> users = userRepository.getUserByUsername(request.username());
        if(users.size()>0){
            User u = users.get(0);
            Long id = u.getUserId();
            //List<Book> books = bookRepository.
        }else{
            ResponseEntity.notFound().build();
        }
        return
        //bookRepository.getBooksBy
    }
    */
}
