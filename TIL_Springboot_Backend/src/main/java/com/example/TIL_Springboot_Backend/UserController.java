package com.example.TIL_Springboot_Backend;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpRequest;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;
import com.example.TIL_Springboot_Backend.DTOs.*;

import java.util.List;

@RestController
public class UserController {



    @Autowired
    private UserRepository userRepository;

    @PostMapping("api/signup")
    public ResponseEntity<String> signup(@RequestBody UserRequest request){
        User u = new User(request.username(), request.password());
        userRepository.save(u);
        System.out.printf("\n\nusername: %s  ,  password: %s", u.getUsername(), u.getPassword());
        return ResponseEntity.ok("User created");
    }

    @PostMapping("api/login")
    public ResponseEntity<User> login(@RequestBody UserRequest request){
        List<User> users = userRepository.getUserByUsername(request.username());

        if(users.size()>0){
            User u = users.get(0);
            System.out.printf("Fetched user: %s", u.getUsername());
            return ResponseEntity.ok(u);
        }else{
            System.out.printf("User with username %s not found", request.username());
            return ResponseEntity.notFound().build();
        }
    }
}